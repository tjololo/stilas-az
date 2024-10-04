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
	apim "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ProductApiVersionSpec defines the desired state of ProductApiVersion
type ProductApiVersionSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//VName - Version identifier for the API Version.
	//+kubebuilder:validation:Required
	Name string `json:"foo,omitempty"`
	//VersionDescription - Description of the API Version.
	//+kubebuilder:validation:Optional
	Description *string `json:"versionDescription,omitempty"`
	//VersioningScheme - An value that determines where the API Version identifer will be located in a HTTP request.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:="Segment"
	//+kubebuilder:validation:Enum:=Header;Query;Segment
	VersioningScheme APIVersionScheme `json:"versioningScheme,omitempty"`
}

// ProductApiVersionStatus defines the observed state of ProductApiVersion
type ProductApiVersionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//+kubebuilder:validation:Optional
	ApiVersionSetID string `json:"apiVersionSetID,omitempty"`
	//+kubebuilder:validation:Optional
	//ProvisioningState - The provisioning state of the API. Possible values are: Creating, Succeeded, Failed, Updating, Deleting, and Deleted.
	//+kubebuilder:validation:Optional
	ProvisioningState string `json:"provisioningState,omitempty"`
	//ResumeToken - The token used to track long-running operations.
	//+kubebuilder:validation:Optional
	ResumeToken string `json:"pollerToken,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ProductApiVersion is the Schema for the productapiversions API
type ProductApiVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProductApiVersionSpec   `json:"spec,omitempty"`
	Status ProductApiVersionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProductApiVersionList contains a list of ProductApiVersion
type ProductApiVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProductApiVersion `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProductApiVersion{}, &ProductApiVersionList{})
}

type APIVersionScheme string

const (
	// APIVersionSetContractDetailsVersioningSchemeHeader - The API Version is passed in a HTTP header.
	APIVersionSetContractDetailsVersioningSchemeHeader APIVersionScheme = "Header"
	// APIVersionSetContractDetailsVersioningSchemeQuery - The API Version is passed in a query parameter.
	APIVersionSetContractDetailsVersioningSchemeQuery APIVersionScheme = "Query"
	// APIVersionSetContractDetailsVersioningSchemeSegment - The API Version is passed in a path segment.
	APIVersionSetContractDetailsVersioningSchemeSegment APIVersionScheme = "Segment"
)

func (a *APIVersionScheme) AzureAPIVersionScheme() *apim.VersioningScheme {
	if a == nil {
		return nil
	}
	apiVersionScheme := apim.VersioningScheme(*a)
	return &apiVersionScheme
}

func (a *APIVersionScheme) AzureAPIVersionSetContractDetailsVersioningScheme() *apim.APIVersionSetContractDetailsVersioningScheme {
	if a == nil {
		return nil
	}
	apiVersionScheme := apim.APIVersionSetContractDetailsVersioningScheme(*a)
	return &apiVersionScheme
}
