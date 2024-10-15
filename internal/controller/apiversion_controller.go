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
		logger.Error(err, "unable to fetch ApiVersion")
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
	_, err = r.apimClient.GetApi(ctx, getApiVersionName(apiVersion), nil)
	if azure.IgnoreNotFound(err) != nil {
		logger.Error(err, "Failed to get API")
		return ctrl.Result{}, err
	} else {
		latestSha, shaErr := utils.Sha256FromUrlContent(*apiVersion.Spec.Content)
		if shaErr != nil {
			logger.Error(err, "Failed to get content sha")
			return ctrl.Result{}, err
		}
		if apiVersion.Status.LastAppliedSpecSha != latestSha || azure.IsNotFoundError(err) {
			return r.createUpdateApimApi(ctx, apiVersion)
		}
		logger.Info("No changes detected")
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
	apimApiParams := apim.APICreateOrUpdateParameter{
		Properties: &apim.APICreateOrUpdateProperties{
			Path:                 &apiVesrion.Spec.Path,
			APIType:              apiVesrion.Spec.APIType.AzureApiType(),
			Contact:              apiVesrion.Spec.Contact.AzureAPIContactInformation(),
			Description:          &apiVesrion.Spec.Description,
			DisplayName:          &apiVesrion.Spec.DisplayName,
			Format:               apiVesrion.Spec.ContentFormat.AzureContentFormat(),
			IsCurrent:            toPointer(true),
			Protocols:            []*apim.Protocol{toPointer(apim.ProtocolHTTPS)},
			ServiceURL:           &apiVesrion.Spec.ServiceUrl,
			SubscriptionRequired: apiVesrion.Spec.SubscriptionRequired,
			Value:                apiVesrion.Spec.Content,
			APIVersionSetID:      toPointer(apiVesrion.Spec.ApiVersionSetId),
			APIVersion:           apiVesrion.Spec.Name,
		},
	}
	poller, err := r.apimClient.CreateUpdateApi(
		ctx,
		getApiVersionName(apiVesrion),
		apimApiParams,
		&apim.APIClientBeginCreateOrUpdateOptions{ResumeToken: resumeToken})

	logger.Info("Watching LR operation")
	status, _, token, err := azure.StartResumeOperation(ctx, poller)

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
		if *apiVesrion.Spec.ContentFormat == apimv1alpha1.ContentFormatOpenapiJSONLink {
			apiVesrion.Status.LastAppliedSpecSha, err = utils.Sha256FromUrlContent(*apiVesrion.Spec.Content)
		}
		err = r.Status().Update(ctx, &apiVesrion)
		if err != nil {
			logger.Error(err, "Failed to update status")
			return ctrl.Result{}, err
		}
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}
	return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
}
