//go:build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIContactInformation) DeepCopyInto(out *APIContactInformation) {
	*out = *in
	if in.Email != nil {
		in, out := &in.Email, &out.Email
		*out = new(string)
		**out = **in
	}
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.URL != nil {
		in, out := &in.URL, &out.URL
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIContactInformation.
func (in *APIContactInformation) DeepCopy() *APIContactInformation {
	if in == nil {
		return nil
	}
	out := new(APIContactInformation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Api) DeepCopyInto(out *Api) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Api.
func (in *Api) DeepCopy() *Api {
	if in == nil {
		return nil
	}
	out := new(Api)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Api) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApiList) DeepCopyInto(out *ApiList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Api, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApiList.
func (in *ApiList) DeepCopy() *ApiList {
	if in == nil {
		return nil
	}
	out := new(ApiList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ApiList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApiPolicySpec) DeepCopyInto(out *ApiPolicySpec) {
	*out = *in
	if in.PolicyContent != nil {
		in, out := &in.PolicyContent, &out.PolicyContent
		*out = new(string)
		**out = **in
	}
	if in.PolicyFormat != nil {
		in, out := &in.PolicyFormat, &out.PolicyFormat
		*out = new(PolicyFormat)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApiPolicySpec.
func (in *ApiPolicySpec) DeepCopy() *ApiPolicySpec {
	if in == nil {
		return nil
	}
	out := new(ApiPolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApiSpec) DeepCopyInto(out *ApiSpec) {
	*out = *in
	if in.Description != nil {
		in, out := &in.Description, &out.Description
		*out = new(string)
		**out = **in
	}
	if in.ApiType != nil {
		in, out := &in.ApiType, &out.ApiType
		*out = new(APIType)
		**out = **in
	}
	if in.Contact != nil {
		in, out := &in.Contact, &out.Contact
		*out = new(APIContactInformation)
		(*in).DeepCopyInto(*out)
	}
	if in.Versions != nil {
		in, out := &in.Versions, &out.Versions
		*out = make([]ApiVersionSubSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApiSpec.
func (in *ApiSpec) DeepCopy() *ApiSpec {
	if in == nil {
		return nil
	}
	out := new(ApiSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApiStatus) DeepCopyInto(out *ApiStatus) {
	*out = *in
	if in.VersionStates != nil {
		in, out := &in.VersionStates, &out.VersionStates
		*out = make(map[string]ApiVersionStatus, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApiStatus.
func (in *ApiStatus) DeepCopy() *ApiStatus {
	if in == nil {
		return nil
	}
	out := new(ApiStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApiVersion) DeepCopyInto(out *ApiVersion) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApiVersion.
func (in *ApiVersion) DeepCopy() *ApiVersion {
	if in == nil {
		return nil
	}
	out := new(ApiVersion)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ApiVersion) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApiVersionList) DeepCopyInto(out *ApiVersionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ApiVersion, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApiVersionList.
func (in *ApiVersionList) DeepCopy() *ApiVersionList {
	if in == nil {
		return nil
	}
	out := new(ApiVersionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ApiVersionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApiVersionSpec) DeepCopyInto(out *ApiVersionSpec) {
	*out = *in
	if in.APIType != nil {
		in, out := &in.APIType, &out.APIType
		*out = new(APIType)
		**out = **in
	}
	if in.Contact != nil {
		in, out := &in.Contact, &out.Contact
		*out = new(APIContactInformation)
		(*in).DeepCopyInto(*out)
	}
	in.ApiVersionSubSpec.DeepCopyInto(&out.ApiVersionSubSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApiVersionSpec.
func (in *ApiVersionSpec) DeepCopy() *ApiVersionSpec {
	if in == nil {
		return nil
	}
	out := new(ApiVersionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApiVersionStatus) DeepCopyInto(out *ApiVersionStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApiVersionStatus.
func (in *ApiVersionStatus) DeepCopy() *ApiVersionStatus {
	if in == nil {
		return nil
	}
	out := new(ApiVersionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApiVersionSubSpec) DeepCopyInto(out *ApiVersionSubSpec) {
	*out = *in
	if in.Name != nil {
		in, out := &in.Name, &out.Name
		*out = new(string)
		**out = **in
	}
	if in.ServiceUrl != nil {
		in, out := &in.ServiceUrl, &out.ServiceUrl
		*out = new(string)
		**out = **in
	}
	if in.Products != nil {
		in, out := &in.Products, &out.Products
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ContentFormat != nil {
		in, out := &in.ContentFormat, &out.ContentFormat
		*out = new(ContentFormat)
		**out = **in
	}
	if in.Content != nil {
		in, out := &in.Content, &out.Content
		*out = new(string)
		**out = **in
	}
	if in.SubscriptionRequired != nil {
		in, out := &in.SubscriptionRequired, &out.SubscriptionRequired
		*out = new(bool)
		**out = **in
	}
	if in.Protocols != nil {
		in, out := &in.Protocols, &out.Protocols
		*out = make([]Protocol, len(*in))
		copy(*out, *in)
	}
	if in.IsCurrent != nil {
		in, out := &in.IsCurrent, &out.IsCurrent
		*out = new(bool)
		**out = **in
	}
	if in.Policy != nil {
		in, out := &in.Policy, &out.Policy
		*out = new(ApiPolicySpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApiVersionSubSpec.
func (in *ApiVersionSubSpec) DeepCopy() *ApiVersionSubSpec {
	if in == nil {
		return nil
	}
	out := new(ApiVersionSubSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Backend) DeepCopyInto(out *Backend) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Backend.
func (in *Backend) DeepCopy() *Backend {
	if in == nil {
		return nil
	}
	out := new(Backend)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Backend) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackendList) DeepCopyInto(out *BackendList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Backend, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackendList.
func (in *BackendList) DeepCopy() *BackendList {
	if in == nil {
		return nil
	}
	out := new(BackendList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BackendList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackendSpec) DeepCopyInto(out *BackendSpec) {
	*out = *in
	if in.Description != nil {
		in, out := &in.Description, &out.Description
		*out = new(string)
		**out = **in
	}
	if in.ValidateCertificateChain != nil {
		in, out := &in.ValidateCertificateChain, &out.ValidateCertificateChain
		*out = new(bool)
		**out = **in
	}
	if in.ValidateCertificateName != nil {
		in, out := &in.ValidateCertificateName, &out.ValidateCertificateName
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackendSpec.
func (in *BackendSpec) DeepCopy() *BackendSpec {
	if in == nil {
		return nil
	}
	out := new(BackendSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BackendStatus) DeepCopyInto(out *BackendStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BackendStatus.
func (in *BackendStatus) DeepCopy() *BackendStatus {
	if in == nil {
		return nil
	}
	out := new(BackendStatus)
	in.DeepCopyInto(out)
	return out
}
