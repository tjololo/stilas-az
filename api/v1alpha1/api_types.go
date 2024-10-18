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

// ApiSpec defines the desired state of Api
type ApiSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//DisplayName - The display name of the API. This name is used by the developer portal as the API name.
	//+kubebuilder:validation:Required
	DisplayName string `json:"displayName,omitempty"`
	//Description - Description of the API. May include its purpose, where to get more information, and other relevant information.
	//+kubebuilder:validation:Optional
	Description *string `json:"description,omitempty"`
	//VersioningScheme - Indicates the versioning scheme used for the API. Possible values include, but are not limited to, "Segment", "Query", "Header". Default value is "Segment".
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:="Segment"
	//+kubebuilder:validation:Enum:=Header;Query;Segment
	VersioningScheme APIVersionScheme `json:"versioningScheme,omitempty"`
	//Path - API prefix. The value is combined with the API version to form the URL of the API endpoint.
	//+kubebuilder:validation:Required
	Path string `json:"path,omitempty"`
	//ApiType - Type of API.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:="http"
	//+default:value:"http"
	//+kubebuilder:validation:Enum:=graphql;http;websocket
	ApiType *APIType `json:"apiType,omitempty"`
	//Contact - Contact details of the API owner.
	//+kubebuilder:validation:Optional
	Contact *APIContactInformation `json:"contact,omitempty"`
	//Versions - A list of API versions associated with the API. If the API is specified using the OpenAPI definition, then the API version is set by the version field of the OpenAPI definition.
	//+kubebuilder:validation:Required
	Versions []ApiVersionSubSpec `json:"versions,omitempty"`
}

// ApiStatus defines the observed state of Api
type ApiStatus struct {
	//ProvisioningState - The provisioning state of the API. Possible values are: Creating, Succeeded, Failed, Updating, Deleting, and Deleted.
	//+kubebuilder:validation:Optional
	ProvisioningState string `json:"provisioningState,omitempty"`
	//+kubebuilder:validation:Optional
	ApiVersionSetID string `json:"apiVersionSetID,omitempty"`
	//VersionStates - A list of API Version deployed in the API Management service.
	//+kubebuilder:validation:Optional
	VersionStates map[string]ApiVersionStatus `json:"versionStates,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Api is the Schema for the apis API
type Api struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiSpec   `json:"spec,omitempty"`
	Status ApiStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ApiList contains a list of Api
type ApiList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Api `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Api{}, &ApiList{})
}
