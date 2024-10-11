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
	apimv1alpha1 "github.com/tjololo/stilas-az/api/v1alpha1"
	"github.com/tjololo/stilas-az/internal/azure"
	"github.com/tjololo/stilas-az/internal/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
)

// ProductApiReconciler reconciles a ProductApi object
type ProductApiReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	apimClient *azure.APIMClient
}

// +kubebuilder:rbac:groups=apim.azure.stilas.418.cloud,resources=productapis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apim.azure.stilas.418.cloud,resources=productapis/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apim.azure.stilas.418.cloud,resources=productapis/finalizers,verbs=update
// +kubebuilder:rbac:groups=apim.azure.stilas.418.cloud,resources=productapiversions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apim.azure.stilas.418.cloud,resources=productapiversions/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ProductApi object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *ProductApiReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var productApi apimv1alpha1.ProductApi
	if err := r.Get(ctx, req.NamespacedName, &productApi); err != nil {
		logger.Error(err, "unable to fetch CronJob")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	subscriptionID, resourcesGroup, apimName, err := getConfigFromEnv()
	if err != nil {
		logger.Error(err, "Failed to get configuration. No reason to requeue")
		return ctrl.Result{}, nil
	}
	r.apimClient, err = azure.NewAPIMClient(azure.ApimClientConfig{
		SubscriptionId:  subscriptionID,
		ResourceGroup:   resourcesGroup,
		ApimServiceName: apimName,
	})
	if err != nil {
		logger.Error(err, "Failed to create APIM client")
		return ctrl.Result{}, err
	}

	if productApi.Spec.ApiVersion != "" {
		var productApiVersion apimv1alpha1.ProductApiVersion
		if err := r.Get(ctx, client.ObjectKey{Namespace: productApi.Namespace, Name: productApi.Name}, &productApiVersion); err != nil {
			if client.IgnoreNotFound(err) != nil {
				logger.Error(err, "Failed to get product api version")
				return ctrl.Result{}, err
			}
			return r.createApimApiVersionSet(ctx, productApi)
		} else {
			if productApiVersion.Status.ProvisioningState != "Succeeded" {
				logger.Info("Product API Version not yet provisioned. Requeuing")
				return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
			}
			if productApi.Annotations == nil || productApi.Annotations["apim.azure.stilas.418.cloud/productapiversion"] != productApiVersion.Status.ApiVersionSetID {
				productApi.Annotations["apim.azure.stilas.418.cloud/productapiversion"] = productApiVersion.Status.ApiVersionSetID
				if err := r.Update(ctx, &productApi); err != nil {
					logger.Error(err, "Failed to set annotation of product api")
					return ctrl.Result{}, err
				}
			}
		}
	}

	_, err = r.apimClient.GetApi(ctx, fmt.Sprintf("%s-%s", productApi.Namespace, productApi.Name), nil)
	if azure.IgnoreNotFound(err) != nil {
		logger.Error(err, "Failed to get API")
		return ctrl.Result{}, err
	} else if azure.IsNotFoundError(err) {
		return r.createUpdateApimApi(ctx, productApi)
	}
	specSha, err := utils.Sha256FromUrlContent(*productApi.Spec.Content)
	if err != nil {
		logger.Error(err, "Failed to get content sha")
		return ctrl.Result{}, err
	}
	if specSha != productApi.Status.LastAppliedSpecSha {
		return r.createUpdateApimApi(ctx, productApi)
	}
	logger.Info("No changes detected. Requeuing after 10 minutes")
	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func (r *ProductApiReconciler) createApimApiVersionSet(ctx context.Context, productApi apimv1alpha1.ProductApi) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	productApiVersion := apimv1alpha1.ProductApiVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      productApi.Name,
			Namespace: productApi.Namespace,
		},
		Spec: apimv1alpha1.ProductApiVersionSpec{
			Name:             productApi.Spec.DisplayName,
			VersioningScheme: *productApi.Spec.ApiVersionScheme,
		},
	}
	if err := controllerutil.SetControllerReference(&productApi, &productApiVersion, r.Scheme); err != nil {
		logger.Error(err, "Failed to set controller reference")
		return ctrl.Result{}, err
	}
	if err := r.Create(ctx, &productApiVersion); err != nil {
		logger.Error(err, "Failed to create product api version")
		return ctrl.Result{}, err
	}
	productApi.Status.ProvisioningState = "Provisioning"
	if err := r.Status().Update(ctx, &productApi); err != nil {
		logger.Error(err, "Failed to update status")
		return ctrl.Result{}, err
	}
	logger.Info("Product API Version created. Requeuing")
	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

func (r *ProductApiReconciler) createUpdateApimApi(ctx context.Context, productApi apimv1alpha1.ProductApi) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	resumeToken := productApi.Status.ResumeToken
	logger.Info("Creating or updating API")
	apimApiParams := apim.APICreateOrUpdateParameter{
		Properties: &apim.APICreateOrUpdateProperties{
			Path:                 &productApi.Spec.Path,
			APIType:              productApi.Spec.ApiType.AzureApiType(),
			Contact:              productApi.Spec.Contact.AzureAPIContactInformation(),
			Description:          &productApi.Spec.Description,
			DisplayName:          &productApi.Spec.DisplayName,
			Format:               productApi.Spec.ContentFormat.AzureContentFormat(),
			IsCurrent:            toPointer(true),
			Protocols:            []*apim.Protocol{toPointer(apim.ProtocolHTTPS)},
			ServiceURL:           &productApi.Spec.ServiceURL,
			SubscriptionRequired: productApi.Spec.SubscriptionRequired,
			Value:                productApi.Spec.Content,
		},
	}
	if productApi.Annotations != nil && productApi.Annotations["apim.azure.stilas.418.cloud/productapiversion"] != "" {
		apimApiParams.Properties.APIVersionSetID = toPointer(productApi.Annotations["apim.azure.stilas.418.cloud/productapiversion"])
		apimApiParams.Properties.APIVersion = &productApi.Spec.ApiVersion
	}
	poller, err := r.apimClient.CreateUpdateApi(
		ctx,
		fmt.Sprintf("%s-%s", productApi.Namespace, productApi.Name),
		apimApiParams,
		&apim.APIClientBeginCreateOrUpdateOptions{ResumeToken: resumeToken})

	logger.Info("Watching LR operation")
	done, _, token, err := azure.StartResumeOperation(ctx, poller)
	if err != nil {
		logger.Error(err, "Failed to watch LR operation")
		productApi.Status.ResumeToken = ""
		productApi.Status.ProvisioningState = "Failed"
		err = r.Status().Update(ctx, &productApi)
		if err != nil {
			logger.Error(err, "Failed to update status")
		}
		return ctrl.Result{}, err
	}

	if done {
		logger.Info("Operation completed")
		productApi.Status.ResumeToken = ""
		productApi.Status.ProvisioningState = "Succeeded"
		if *productApi.Spec.ContentFormat == apimv1alpha1.ContentFormatOpenapiJSONLink {
			productApi.Status.LastAppliedSpecSha, err = utils.Sha256FromUrlContent(*productApi.Spec.Content)
		}
		err = r.Status().Update(ctx, &productApi)
		if err != nil {
			logger.Error(err, "Failed to update status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	} else {
		productApi.Status.ProvisioningState = "Provisioning"
		productApi.Status.ResumeToken = token
		err = r.Status().Update(ctx, &productApi)
		if err != nil {
			logger.Error(err, "Failed to update status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProductApiReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apimv1alpha1.ProductApi{}).
		Owns(&apimv1alpha1.ProductApiVersion{}).
		Complete(r)
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
