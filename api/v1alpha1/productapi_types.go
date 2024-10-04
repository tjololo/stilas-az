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

// ProductApiSpec defines the desired state of ProductApi
type ProductApiSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//DisplayName - The display name of the API. This name is used by the developer portal as the API name.
	//+kubebuilder:validation:Required
	DisplayName string `json:"displayName,omitempty"`
	//Description - Description of the API. May include its purpose, where to get more information, and other relevant information.
	//+kubebuilder:validation:Required
	Description string `json:"description,omitempty"`
	//ServiceURL - Absolute URL of the backend service implementing this API. Cannot be more than 2000 characters long.
	//+kubebuilder:validation:Required
	ServiceURL string `json:"serviceURL,omitempty"`
	//Path - API prefix. The value is combined with the API version to form the URL of the API endpoint.
	//+kubebuilder:validation:Required
	Path string `json:"path,omitempty"`
	//Products - Products that the API is associated with. Products are groups of APIs.
	//+kubebuilder:validation:Optional
	Products []string `json:"products,omitempty"`
	//Protocols - Describes protocols used by the API. Default value is [https].
	//+kubebuilder:validation:Optional
	ApiVersion string `json:"apiVersion,omitempty"`
	//ApiType - Type of API.
	//+kubebuilder:validation:Optional
	//+kubebuilder:default:="http"
	//+default:value:"http"
	//+kubebuilder:validation:Enum:=graphql;http;websocket
	ApiType *APIType `json:"apiType,omitempty"`
	//Contact - Contact details of the API owner.
	//+kubebuilder:validation:Optional
	Contact *APIContactInformation `json:"contact,omitempty"`
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
}

// ProductApiStatus defines the observed state of ProductApi
type ProductApiStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//ProvisioningState - The provisioning state of the API. Possible values are: Creating, Succeeded, Failed, Updating, Deleting, and Deleted.
	//+kubebuilder:validation:Optional
	ProvisioningState string `json:"provisioningState,omitempty"`
	//ResumeToken - The token used to track long-running operations.
	//+kubebuilder:validation:Optional
	ResumeToken string `json:"pollerToken,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ProductApi is the Schema for the productapis API
type ProductApi struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProductApiSpec   `json:"spec,omitempty"`
	Status ProductApiStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProductApiList contains a list of ProductApi
type ProductApiList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProductApi `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProductApi{}, &ProductApiList{})
}

// APIType - Type of API.
type APIType string

const (
	APITypeGraphql   APIType = "graphql"
	APITypeHTTP      APIType = "http"
	APITypeWebsocket APIType = "websocket"
)

func (a APIType) AzureApiType() *apim.APIType {
	apiType := apim.APIType(a)
	return &apiType
}

// ContentFormat - Format of the Content in which the API is getting imported.
type ContentFormat string

const (
	// ContentFormatGraphqlLink - The GraphQL API endpoint hosted on a publicly accessible internet address.
	ContentFormatGraphqlLink ContentFormat = "graphql-link"
	// ContentFormatOpenapi - The contents are inline and Content Type is a OpenAPI 3.0 YAML Document.
	ContentFormatOpenapi ContentFormat = "openapi"
	// ContentFormatOpenapiJSON - The contents are inline and Content Type is a OpenAPI 3.0 JSON Document.
	ContentFormatOpenapiJSON ContentFormat = "openapi+json"
	// ContentFormatOpenapiJSONLink - The OpenAPI 3.0 JSON document is hosted on a publicly accessible internet address.
	ContentFormatOpenapiJSONLink ContentFormat = "openapi+json-link"
	// ContentFormatOpenapiLink - The OpenAPI 3.0 YAML document is hosted on a publicly accessible internet address.
	ContentFormatOpenapiLink ContentFormat = "openapi-link"
	// ContentFormatSwaggerJSON - The contents are inline and Content Type is a OpenAPI 2.0 JSON Document.
	ContentFormatSwaggerJSON ContentFormat = "swagger-json"
	// ContentFormatSwaggerLinkJSON - The OpenAPI 2.0 JSON document is hosted on a publicly accessible internet address.
	ContentFormatSwaggerLinkJSON ContentFormat = "swagger-link-json"
	// ContentFormatWadlLinkJSON - The WADL document is hosted on a publicly accessible internet address.
	ContentFormatWadlLinkJSON ContentFormat = "wadl-link-json"
	// ContentFormatWadlXML - The contents are inline and Content type is a WADL document.
	ContentFormatWadlXML ContentFormat = "wadl-xml"
)

func (c ContentFormat) AzureContentFormat() *apim.ContentFormat {
	contentFormat := apim.ContentFormat(c)
	return &contentFormat
}

type APIContactInformation struct {
	// The email address of the contact person/organization. MUST be in the format of an email address
	Email *string `json:"email,omitempty"`

	// The identifying name of the contact person/organization
	Name *string `json:"name,omitempty"`

	// The URL pointing to the contact information. MUST be in the format of a URL
	URL *string `json:"url,omitempty"`
}

func (a *APIContactInformation) AzureAPIContactInformation() *apim.APIContactInformation {
	if a == nil {
		return nil
	}
	return &apim.APIContactInformation{
		Email: a.Email,
		Name:  a.Name,
		URL:   a.URL,
	}
}
