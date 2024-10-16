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

// ApiVersionSpec defines the desired state of ApiVersion
type ApiVersionSpec struct {
	ApiVersionSetId   string                 `json:"apiVersionSetId,omitempty"`
	ApiVersionScheme  APIVersionScheme       `json:"apiVersionScheme,omitempty"`
	Path              string                 `json:"path,omitempty"`
	APIType           *APIType               `json:"apiType,omitempty"`
	Contact           *APIContactInformation `json:"contact,omitempty"`
	ApiVersionSubSpec `json:",inline"`
}

// ApiVersionSubSpec defines the desired state of ApiVersion
type ApiVersionSubSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//+kubebuilder:validation:Optional
	Name *string `json:"name,omitempty"`
	//DisplayName - The display name of the API Version. This name is used by the developer portal as the API Version name.
	//+kubebuilder:validation:Required
	DisplayName string `json:"displayName,omitempty"`
	//Description - Description of the API Version. May include its purpose, where to get more information, and other relevant information.
	//+kubebuilder:validation:Optional
	Description string `json:"description,omitempty"`
	//ServiceUrl - Absolute URL of the backend service implementing this API. Cannot be more than 2000 characters long.
	//+kubebuilder:validation:Optional
	ServiceUrl *string `json:"serviceUrl,omitempty"`
	//Products - Products that the API is associated with. Products are groups of APIs.
	//+kubebuilder:validation:Optional
	Products []string `json:"products,omitempty"`
	//ContentFormat - Format of the Content in which the API is getting imported.
	//+kubebuilder:validation:Required
	//+kubebuilder:default:=openapi+json
	ContentFormat *ContentFormat `json:"contentFormat,omitempty"`
	//Content - The contents of the API. The value is a string containing the content of the API.
	//+kubebuilder:validation:Required
	Content *string `json:"content,omitempty"`
	//SubscriptionRquired - Indicates if subscription is required to access the API. Default value is true.
	//+kubebuilder:validation:Required
	//+kubebuilder:default:=true
	SubscriptionRequired *bool `json:"subscriptionRequired,omitempty"`
	//Protocols - Describes protocols over which API is made available.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:={https}
	Protocols []Protocol `json:"protocols,omitempty"`
	//IsCurrent - Indicates if API Version is the current api version.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=true
	IsCurrent *bool `json:"isCurrent,omitempty"`
	//Policy - The API Version Policy description.
	//+kubebuilder:validation:Optional
	Policy *ApiPolicySpec `json:"policies,omitempty"`
}

// ApiPolicySpec defines the desired state of ApiVersion
type ApiPolicySpec struct {
	//PolicyContent - The contents of the Policy as string.
	//+kubebuilder:validation:Required
	PolicyContent *string `json:"policyContent,omitempty"`
	//PolicyFormat - Format of the Policy in which the API is getting imported.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:=xml
	//+kubebuilder:validation:Enum:=xml;xml-link;rawxml;rawxml-link
	PolicyFormat *PolicyFormat `json:"policyFormat,omitempty"`
}

// ApiVersionStatus defines the observed state of ApiVersion
type ApiVersionStatus struct {
	//ProvisioningState - The provisioning state of the API. Possible values are: Creating, Succeeded, Failed, Updating, Deleting, and Deleted.
	//+kubebuilder:validation:Optional
	ProvisioningState string `json:"provisioningState,omitempty"`
	//ResumeToken - The token used to track long-running operations.
	//+kubebuilder:validation:Optional
	ResumeToken string `json:"pollerToken,omitempty"`
	//LastAppliedSpecSha - The sha256 of the last applied spec.
	//+kubebuilder:validation:Optional
	LastAppliedSpecSha string `json:"lastAppliedSpecSha,omitempty"`
	//LastAppliedPolicySha - The sha256 of the last applied policy.
	//+kubebuilder:validation:Optional
	LastAppliedPolicySha string `json:"lastAppliedPolicySha,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ApiVersion is the Schema for the apiversions API
type ApiVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiVersionSpec   `json:"spec,omitempty"`
	Status ApiVersionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ApiVersionList contains a list of ApiVersion
type ApiVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApiVersion `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApiVersion{}, &ApiVersionList{})
}
