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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apimv1alpha1 "github.com/tjololo/stilas-az/api/v1alpha1"
)

// BackendReconciler reconciles a Backend object
type BackendReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	apimClient *azure.APIMClient
}

// +kubebuilder:rbac:groups=apim.azure.stilas.418.cloud,resources=backends,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apim.azure.stilas.418.cloud,resources=backends/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apim.azure.stilas.418.cloud,resources=backends/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Backend object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *BackendReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// get the backend
	var backend apimv1alpha1.Backend
	if err := r.Get(ctx, req.NamespacedName, &backend); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if !controllerutil.ContainsFinalizer(&backend, "backend.finalizers.stilas.418.cloud") {
		controllerutil.AddFinalizer(&backend, "backend.finalizers.stilas.418.cloud")
		if err := r.Update(ctx, &backend); err != nil {
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
		return ctrl.Result{}, err
	}
	azureBackend, err := r.apimClient.GetBackend(ctx, getBackendName(backend), nil)
	if err != nil {
		if azure.IsNotFoundError(err) {
			logger.Info("Backend not found in Azure, creating")
			createdBackend, err := r.apimClient.CreateUpdateBackend(ctx, getBackendName(backend), toAzureBackend(&backend), nil)
			if err != nil {
				logger.Error(err, "Failed to create backend")
				backend.Status.ProvisioningState = "Failed"
				if errUpdate := r.Status().Update(ctx, &backend); errUpdate != nil {
					logger.Error(err, "Failed to update status")
				}
				return ctrl.Result{}, err
			}
			backend.Status.BackendID = *createdBackend.ID
			backend.Status.ProvisioningState = "Succeeded"
			if errUpdate := r.Status().Update(ctx, &backend); errUpdate != nil {
				logger.Error(err, "Failed to update status")
			}
			return ctrl.Result{}, nil
		} else {
			logger.Error(err, "Failed to get backend")
			return ctrl.Result{}, err
		}
	}
	if backend.DeletionTimestamp != nil {
		logger.Info("Deleting backend")
		_, err := r.apimClient.DeleteBackend(ctx, getBackendName(backend), *azureBackend.ETag, nil)
		if err != nil {
			logger.Error(err, "Failed to delete backend")
			return ctrl.Result{}, err
		}
		controllerutil.RemoveFinalizer(&backend, "backend.finalizers.stilas.418.cloud")
		if err := r.Update(ctx, &backend); err != nil {
			logger.Error(err, "Failed to remove finalizer")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}
	if *azureBackend.Properties.URL != backend.Spec.Url {
		logger.Info("Updating backend")
		updatedBackend, err := r.apimClient.CreateUpdateBackend(ctx, getBackendName(backend), toAzureBackend(&backend), nil)
		if err != nil {
			logger.Error(err, "Failed to update backend")
			backend.Status.ProvisioningState = "Failed"
			if errUpdate := r.Status().Update(ctx, &backend); errUpdate != nil {
				logger.Error(err, "Failed to update status")
			}
			return ctrl.Result{}, err
		}
		backend.Status.BackendID = *updatedBackend.ID
		backend.Status.ProvisioningState = "Succeeded"
		if errUpdate := r.Status().Update(ctx, &backend); errUpdate != nil {
			logger.Error(err, "Failed to update status")
		}
	}
	return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *BackendReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apimv1alpha1.Backend{}).
		Complete(r)
}

func getBackendName(backend apimv1alpha1.Backend) string {
	return fmt.Sprintf("%s-%s", backend.Namespace, backend.Name)
}

func toAzureBackend(backend *apimv1alpha1.Backend) apim.BackendContract {
	return apim.BackendContract{
		Properties: &apim.BackendContractProperties{
			Protocol:    toPointer(apim.BackendProtocolHTTP),
			URL:         toPointer(backend.Spec.Url),
			Description: backend.Spec.Description,
			TLS: &apim.BackendTLSProperties{
				ValidateCertificateChain: backend.Spec.ValidateCertificateChain,
				ValidateCertificateName:  backend.Spec.ValidateCertificateName,
			},
			Title: toPointer(backend.Spec.Title),
		},
	}
}
