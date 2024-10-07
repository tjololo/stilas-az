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
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	apim "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v2"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apimv1alpha1 "github.com/tjololo/stilas-az/api/v1alpha1"
)

// ProductApiVersionReconciler reconciles a ProductApiVersion object
type ProductApiVersionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apim.azure.stilas.418.cloud,resources=productapiversions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apim.azure.stilas.418.cloud,resources=productapiversions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apim.azure.stilas.418.cloud,resources=productapiversions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ProductApiVersion object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *ProductApiVersionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var productApiVersion apimv1alpha1.ProductApiVersion
	if err := r.Get(ctx, req.NamespacedName, &productApiVersion); err != nil {
		logger.Error(err, "unable to fetch CronJob")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if productApiVersion.Status.ProvisioningState == "Succeeded" {
		logger.Info("Resource already provisioned successfully")
		return ctrl.Result{}, nil
	}
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		logger.Error(err, "Failed to get Azure credentials")
		return ctrl.Result{RequeueAfter: 5 * time.Minute}, err
	}
	subscriptionID, resourcesGroup, apimName, err := getConfigFromEnv()
	if err != nil {
		logger.Error(err, "Failed to get configuration. No reason to requeue")
		return ctrl.Result{}, nil
	}
	apimanagementClientFactory, err := apim.NewClientFactory(subscriptionID, cred, nil)
	if err != nil {
		logger.Error(err, "Failed to create apimanagement client factory")
		return ctrl.Result{RequeueAfter: 5 * time.Minute}, err
	}
	apiVersionClient := apimanagementClientFactory.NewAPIVersionSetClient()

	result, err := apiVersionClient.CreateOrUpdate(
		ctx,
		resourcesGroup,
		apimName,
		productApiVersion.Spec.Name,
		apim.APIVersionSetContract{
			Properties: &apim.APIVersionSetContractProperties{
				DisplayName:      &productApiVersion.Spec.Name,
				VersioningScheme: productApiVersion.Spec.VersioningScheme.AzureAPIVersionScheme(),
				Description:      productApiVersion.Spec.Description,
			},
			Name: &productApiVersion.Name,
		},
		nil)

	if err != nil {
		logger.Error(err, "Failed to create or update API version")
		return ctrl.Result{RequeueAfter: 5 * time.Minute}, err
	}
	if result.ID == nil {
		logger.Info("No result returned")
		return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
	}

	productApiVersion.Status.ProvisioningState = "Succeeded"
	productApiVersion.Status.ApiVersionSetID = *result.ID

	err = r.Status().Update(ctx, &productApiVersion)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ProductApiVersionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apimv1alpha1.ProductApiVersion{}).
		Complete(r)
}
