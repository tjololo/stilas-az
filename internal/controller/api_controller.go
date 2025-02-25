/*
Copyright 2024 tjololo.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	apim "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v2"
	"github.com/tjololo/stilas-az/internal/azure"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apimv1alpha1 "github.com/tjololo/stilas-az/api/v1alpha1"
)

type newApimCLient func(config azure.ApimClientConfig) (*azure.APIMClient, error)

// ApiReconciler reconciles a Api object
type ApiReconciler struct {
	client.Client
	NewClient  newApimCLient
	Scheme     *runtime.Scheme
	apimClient *azure.APIMClient
}

// +kubebuilder:rbac:groups=apim.azure.stilas.418.cloud,resources=apis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apim.azure.stilas.418.cloud,resources=apis/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apim.azure.stilas.418.cloud,resources=apis/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Api object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *ApiReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var api apimv1alpha1.Api
	if err := r.Get(ctx, req.NamespacedName, &api); err != nil {
		if client.IgnoreNotFound(err) != nil {
			logger.Error(err, "unable to fetch Api")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if !controllerutil.ContainsFinalizer(&api, "api.finalizers.stilas.418.cloud") {
		controllerutil.AddFinalizer(&api, "api.finalizers.stilas.418.cloud")
		err := r.Update(ctx, &api)
		if err != nil {
			logger.Error(err, "Failed to add finalizer")
			return ctrl.Result{}, err
		}
	}
	logger.Info("Reconciling Api")
	subscriptionID, resourcesGroup, apimName, err := getConfigFromEnv()
	if err != nil {
		logger.Error(err, "Failed to get configuration. No reason to requeue")
		return ctrl.Result{}, nil
	}
	r.apimClient, err = r.NewClient(azure.ApimClientConfig{
		SubscriptionId:  subscriptionID,
		ResourceGroup:   resourcesGroup,
		ApimServiceName: apimName,
	})
	if err != nil {
		logger.Error(err, "Failed to create APIM client")
		return ctrl.Result{}, err
	}
	apiName := getApiName(&api)
	if api.DeletionTimestamp != nil {
		logger.Info("Deleting API")
		//Get all owned resources and delete them first
		done, err := r.deleteOwnedResources(ctx, &api)
		if err != nil {
			logger.Error(err, "Failed to delete owned resources")
			return ctrl.Result{}, err
		}
		if !done {
			logger.Info("Owned resources not yet deleted")
			return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
		}
		err = r.deleteAzureResources(ctx, apiName)
		if err != nil {
			logger.Error(err, "Failed to delete Azure resources")
			return ctrl.Result{}, err
		}
		controllerutil.RemoveFinalizer(&api, "api.finalizers.stilas.418.cloud")
		err = r.Update(ctx, &api)
		if err != nil {
			logger.Error(err, "Failed to remove finalizer")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	var resId *string

	getRes, err := r.apimClient.GetApiVersionSet(ctx, apiName, nil)
	if azure.IsNotFoundError(err) {
		result, err := r.apimClient.CreateUpdateApiVersionSet(
			ctx,
			apiName,
			apim.APIVersionSetContract{
				Properties: &apim.APIVersionSetContractProperties{
					DisplayName:      &api.Spec.DisplayName,
					VersioningScheme: api.Spec.VersioningScheme.AzureAPIVersionScheme(),
					Description:      api.Spec.Description,
				},
				Name: &apiName,
			},
			nil)

		if err != nil {
			logger.Error(err, "Failed to create or update API version")
			return ctrl.Result{}, err
		}
		resId = result.ID
	} else if err != nil {
		logger.Error(err, "Failed to get API version")
		return ctrl.Result{}, err
	} else {
		resId = getRes.ID
	}
	if resId == nil {
		logger.Info("No result returned")
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	api.Status.ProvisioningState = "Succeeded"
	api.Status.ApiVersionSetID = *resId
	err = r.reconcileVersions(ctx, &api)
	if err != nil {
		logger.Error(err, "Failed to reconcile versions")
		return ctrl.Result{}, err
	}
	err = r.Status().Update(ctx, &api)
	if err != nil {
		logger.Error(err, "Failed to update status of product api version")
		return ctrl.Result{}, err
	}
	return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApiReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Create an index for the ownerReferences.uid field
	if err := mgr.GetFieldIndexer().IndexField(context.TODO(), &apimv1alpha1.ApiVersion{}, "metadata.ownerReferences.uid", func(rawObj client.Object) []string {
		// Extract the owner UID from the ownerReferences
		apiVersion := rawObj.(*apimv1alpha1.ApiVersion)
		ownerRefs := apiVersion.GetOwnerReferences()
		if len(ownerRefs) == 0 {
			return nil
		}
		return []string{string(ownerRefs[0].UID)}
	}); err != nil {
		return err
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&apimv1alpha1.Api{}).
		Owns(&apimv1alpha1.ApiVersion{}).
		Complete(r)
}

func (r *ApiReconciler) reconcileVersions(ctx context.Context, api *apimv1alpha1.Api) error {
	logger := log.FromContext(ctx)
	for _, version := range api.Spec.Versions {
		versionSpecifier := version.Name
		if versionSpecifier == nil || *versionSpecifier == "" {
			versionSpecifier = toPointer("default")
		}
		versionName := fmt.Sprintf("%s-%s", getApiName(api), *versionSpecifier)
		var apiVersion apimv1alpha1.ApiVersion
		if err := r.Get(ctx, client.ObjectKey{Namespace: api.Namespace, Name: versionName}, &apiVersion); err != nil {
			if client.IgnoreNotFound(err) != nil {
				logger.Error(err, "Failed to get product api version")
				return err
			}
			apiVersion = createApiVersionResource(versionName, api, version)
			if err := controllerutil.SetControllerReference(api, &apiVersion, r.Scheme); err != nil {
				logger.Error(err, "Failed to set controller reference")
				return err
			}
			if err := r.Create(ctx, &apiVersion); err != nil {
				logger.Error(err, "Failed to create product api version")
				return err
			}
			continue
		} else {
			if newApi := createApiVersionResource(versionName, api, version); apiVersion.RequireUpdate(newApi) {
				logger.Info("Updating product api version")
				apiVersion.Spec = newApi.Spec
				if err := r.Update(ctx, &apiVersion); err != nil {
					logger.Error(err, "Failed to update product api version")
					return err
				}
			}
			if api.Status.VersionStates == nil {
				api.Status.VersionStates = make(map[string]apimv1alpha1.ApiVersionStatus)
			}
			api.Status.VersionStates[versionName] = apiVersion.Status
		}
	}

	return nil
}

func (r *ApiReconciler) deleteOwnedResources(ctx context.Context, api *apimv1alpha1.Api) (done bool, err error) {
	var versions apimv1alpha1.ApiVersionList
	apiVersionErr := r.List(ctx, &versions, client.InNamespace(api.Namespace), client.MatchingFields{"metadata.ownerReferences.uid": string(api.GetUID())})
	if client.IgnoreNotFound(apiVersionErr) != nil {
		return false, apiVersionErr
	}
	for _, version := range versions.Items {
		if version.DeletionTimestamp != nil {
			deleteErr := r.Delete(ctx, &version)
			if deleteErr != nil {
				return false, deleteErr
			}
		}
	}
	return len(versions.Items) == 0, nil
}

func (r *ApiReconciler) deleteAzureResources(ctx context.Context, apiName string) error {
	_, err := r.apimClient.GetApiVersionSet(ctx, apiName, nil)
	if azure.IgnoreNotFound(err) != nil {
		return fmt.Errorf("failed to get API version set: %w", err)
	}
	if err == nil {
		_, err = r.apimClient.DeleteApiVersionSet(ctx, apiName, "*", nil)
		if err != nil {
			return fmt.Errorf("failed to delete API version set: %w", err)
		}
	}
	return nil
}

func createApiVersionResource(versionName string, api *apimv1alpha1.Api, version apimv1alpha1.ApiVersionSubSpec) apimv1alpha1.ApiVersion {
	return apimv1alpha1.ApiVersion{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      versionName,
			Namespace: api.Namespace,
		},
		Spec: apimv1alpha1.ApiVersionSpec{
			ApiVersionSetId:   api.Status.ApiVersionSetID,
			ApiVersionScheme:  api.Spec.VersioningScheme,
			Path:              api.Spec.Path,
			APIType:           api.Spec.ApiType,
			ApiVersionSubSpec: version,
		},
	}
}

func getApiName(api *apimv1alpha1.Api) string {
	return fmt.Sprintf("%s-%s", api.Namespace, api.Name)
}

func toPointer[T any](t T) *T {
	return &t
}

func getConfigFromEnv() (subscriptionID string, resourcesGroup string, apimName string, err error) {
	subscriptionID = os.Getenv("STILAS_AZ_SUBSCRIPTION_ID")
	if subscriptionID == "" {
		err = fmt.Errorf("STILAS_AZ_SUBSCRIPTION_ID must be set")
		return
	}
	resourcesGroup = os.Getenv("STILAS_AZ_RESOURCE_GROUP")
	if resourcesGroup == "" {
		err = fmt.Errorf("STILAS_AZ_RESOURCE_GROUP must be set")
		return
	}
	apimName = os.Getenv("STILAS_AZ_APIM_NAME")
	if apimName == "" {
		err = fmt.Errorf("AZURE_APIM_NAME must be set")
		return
	}
	return
}
