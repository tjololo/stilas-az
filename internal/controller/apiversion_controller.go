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
	"github.com/tjololo/stilas-az/internal/utils"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apimv1alpha1 "github.com/tjololo/stilas-az/api/v1alpha1"
)

// ApiVersionReconciler reconciles a ApiVersion object
type ApiVersionReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	apimClient *azure.APIMClient
}

// +kubebuilder:rbac:groups=apim.azure.stilas.418.cloud,resources=apiversions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apim.azure.stilas.418.cloud,resources=apiversions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apim.azure.stilas.418.cloud,resources=apiversions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ApiVersion object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *ApiVersionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var apiVersion apimv1alpha1.ApiVersion
	if err := r.Get(ctx, req.NamespacedName, &apiVersion); err != nil {
		if client.IgnoreNotFound(err) != nil {
			logger.Error(err, "unable to fetch ApiVersion")
		}

		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if !controllerutil.ContainsFinalizer(&apiVersion, "apiversion.finalizers.stilas.418.cloud") {
		controllerutil.AddFinalizer(&apiVersion, "apiversion.finalizers.stilas.418.cloud")
		if err := r.Update(ctx, &apiVersion); err != nil {
			logger.Error(err, "Failed to add finalizer")
			return ctrl.Result{}, err
		}
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
	_, err = r.apimClient.GetApi(ctx, getApiVersionName(apiVersion), nil)
	if apiVersion.DeletionTimestamp != nil {
		return r.deleteApiVersion(ctx, apiVersion)
	}
	if azure.IgnoreNotFound(err) != nil {
		logger.Error(err, "Failed to get API")
		return ctrl.Result{}, err
	} else {
		latestSha, shaErr := utils.Sha256FromContent(*apiVersion.Spec.Content)
		if shaErr != nil {
			logger.Error(err, "Failed to get content sha")
			return ctrl.Result{}, err
		}
		if apiVersion.Status.LastAppliedSpecSha != latestSha || azure.IsNotFoundError(err) {
			return r.createUpdateApimApi(ctx, apiVersion)
		}
		_, policyErr := r.apimClient.GetApiPolicy(ctx, getApiVersionName(apiVersion), nil)
		lastPolicySha, shaErr := utils.Sha256FromContent(*apiVersion.Spec.Policy.PolicyContent)
		if shaErr != nil {
			logger.Error(err, "Failed to get policy sha")
			return ctrl.Result{}, err
		}
		if apiVersion.Spec.Policy != nil && apiVersion.Status.LastAppliedPolicySha != lastPolicySha || azure.IsNotFoundError(policyErr) {
			if err := r.createUpdatePolicy(ctx, apiVersion); err != nil {
				logger.Error(err, "Failed to create/update policy")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApiVersionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apimv1alpha1.ApiVersion{}).
		Complete(r)
}

func getApiVersionName(apiVersion apimv1alpha1.ApiVersion) string {
	return fmt.Sprintf("%s-%s", apiVersion.Namespace, apiVersion.Name)
}

func (r *ApiVersionReconciler) createUpdateApimApi(ctx context.Context, apiVesrion apimv1alpha1.ApiVersion) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	resumeToken := apiVesrion.Status.ResumeToken
	logger.Info("Creating or updating API")
	apimApiParams := apiVersionToUpdateParameter(apiVesrion)
	poller, err := r.apimClient.CreateUpdateApi(
		ctx,
		getApiVersionName(apiVesrion),
		apimApiParams,
		&apim.APIClientBeginCreateOrUpdateOptions{ResumeToken: resumeToken})

	if err != nil {
		logger.Error(err, "Failed to create/update API")
		return ctrl.Result{}, err
	}
	logger.Info("Watching LR operation")
	status, _, token, err := azure.StartResumeOperation(ctx, poller)
	if err != nil {
		logger.Error(err, "Failed to watch LR operation")
		return ctrl.Result{}, err
	}

	switch status {
	case azure.OperationStatusFailed:
		logger.Error(err, "Failed to watch LR operation")
		apiVesrion.Status.ResumeToken = ""
		apiVesrion.Status.ProvisioningState = "Failed"
		err = r.Status().Update(ctx, &apiVesrion)
		if err != nil {
			logger.Error(err, "Failed to update status")
		}
		return ctrl.Result{}, err
	case azure.OperationStatusInProgress:
		apiVesrion.Status.ProvisioningState = "Provisioning"
		apiVesrion.Status.ResumeToken = token
		err = r.Status().Update(ctx, &apiVesrion)
		if err != nil {
			logger.Error(err, "Failed to update status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	case azure.OperationStatusSucceeded:
		logger.Info("Operation completed")
		apiVesrion.Status.ResumeToken = ""
		apiVesrion.Status.ProvisioningState = "Succeeded"
		apiVesrion.Status.LastAppliedSpecSha, err = utils.Sha256FromContent(*apiVesrion.Spec.Content)
		err = r.Status().Update(ctx, &apiVesrion)
		if err != nil {
			logger.Error(err, "Failed to update status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
}

func (r *ApiVersionReconciler) createUpdatePolicy(ctx context.Context, apiVersion apimv1alpha1.ApiVersion) error {
	logger := log.FromContext(ctx)
	if apiVersion.Spec.Policy == nil {
		return nil
	}
	logger.Info("Creating or updating policy")
	policy := apiVersion.Spec.Policy
	policyContent := *policy.PolicyContent
	policyFormat := policy.PolicyFormat.AzurePolicyFormat()
	_, err := r.apimClient.CreateUpdateApiPolicy(
		ctx,
		getApiVersionName(apiVersion),
		apim.PolicyContract{
			Properties: &apim.PolicyContractProperties{
				Value:  &policyContent,
				Format: policyFormat,
			}},
		nil,
	)
	if err != nil {
		logger.Error(err, "Failed to create/update policy")
		return err
	}
	apiVersion.Status.LastAppliedPolicySha, err = utils.Sha256FromContent(*apiVersion.Spec.Policy.PolicyContent)
	if err != nil {
		logger.Error(err, "Failed to get policy sha")
		return err
	}
	err = r.Status().Update(ctx, &apiVersion)
	if err != nil {
		logger.Error(err, "Failed to update status")
		return err
	}
	return nil
}

func apiVersionToUpdateParameter(apiVesrion apimv1alpha1.ApiVersion) apim.APICreateOrUpdateParameter {
	return apim.APICreateOrUpdateParameter{
		Properties: &apim.APICreateOrUpdateProperties{
			Path:                 &apiVesrion.Spec.Path,
			APIType:              apiVesrion.Spec.APIType.AzureApiType(),
			Contact:              apiVesrion.Spec.Contact.AzureAPIContactInformation(),
			Description:          &apiVesrion.Spec.Description,
			DisplayName:          &apiVesrion.Spec.DisplayName,
			Format:               apiVesrion.Spec.ContentFormat.AzureContentFormat(),
			IsCurrent:            apiVesrion.Spec.IsCurrent,
			Protocols:            apimv1alpha1.ToApimProtocolSlice(apiVesrion.Spec.Protocols),
			ServiceURL:           apiVesrion.Spec.ServiceUrl,
			SubscriptionRequired: apiVesrion.Spec.SubscriptionRequired,
			Value:                apiVesrion.Spec.Content,
			APIVersionSetID:      toPointer(apiVesrion.Spec.ApiVersionSetId),
			APIVersion:           apiVesrion.Spec.Name,
		},
	}
}

func (r *ApiVersionReconciler) deleteApiVersion(ctx context.Context, apiVersion apimv1alpha1.ApiVersion) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Deleting APIVersion")
	_, err := r.apimClient.DeleteApi(ctx, getApiVersionName(apiVersion), "*", nil)
	if azure.IgnoreNotFound(err) != nil {
		logger.Error(err, "Failed to delete APIVersion")
		return ctrl.Result{}, err
	}
	_, err = r.apimClient.DeleteApiPolicy(ctx, getApiVersionName(apiVersion), "*", nil)
	if azure.IgnoreNotFound(err) != nil {
		logger.Error(err, "Failed to delete policy")
		return ctrl.Result{}, err
	}
	controllerutil.RemoveFinalizer(&apiVersion, "apiversion.finalizers.stilas.418.cloud")
	err = r.Update(ctx, &apiVersion)
	if err != nil {
		logger.Error(err, "Failed to remove finalizer")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
