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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BackendSpec defines the desired state of Backend
type BackendSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//Title - Title of the Backend. May include its purpose, where to get more information, and other relevant information.
	//+kubebuilder:validation:Required
	Title string `json:"title,omitempty"`
	//Description - Description of the Backend. May include its purpose, where to get more information, and other relevant information.
	//+kubebuilder:validation:Optional
	Description *string `json:"description,omitempty"`
	//Url - URL of the Backend.
	//+kubebuilder:validation:Required
	Url string `json:"url,omitempty"`
	//ValidateCertificateChain - Whether to validate the certificate chain when using the backend.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=true
	ValidateCertificateChain *bool `json:"validateCertificateChain,omitempty"`
	//ValidateCertificateName - Whether to validate the certificate name when using the backend.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=true
	ValidateCertificateName *bool `json:"validateCertificateName,omitempty"`
}

// BackendStatus defines the observed state of Backend
type BackendStatus struct {
	//BackendID - The identifier of the Backend.
	//+kubebuilder:validation:Optional
	BackendID string `json:"backendID,omitempty"`
	//ProvisioningState - The provisioning state of the Backend.
	//+kubebuilder:validation:Optional
	ProvisioningState string `json:"provisioningState,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Backend is the Schema for the backends API
type Backend struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BackendSpec   `json:"spec,omitempty"`
	Status BackendStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BackendList contains a list of Backend
type BackendList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Backend `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Backend{}, &BackendList{})
}
