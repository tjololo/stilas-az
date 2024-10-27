package v1alpha1

import (
	"fmt"
	apim "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v2"
	"github.com/tjololo/stilas-az/internal/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a *Api) ToAzureApiVersionSetContract() *apim.APIVersionSetContract {
	return &apim.APIVersionSetContract{
		Properties: &apim.APIVersionSetContractProperties{
			DisplayName:      &a.Spec.DisplayName,
			VersioningScheme: a.Spec.VersioningScheme.AzureAPIVersionScheme(),
			Description:      a.Spec.Description,
		},
		Name: a.GetAzureApiName(),
	}
}

func (a *Api) GetAzureApiName() *string {
	name := fmt.Sprintf("%s-%s", a.Namespace, a.Name)
	return &name
}

func (a *Api) ToApiVersions() []*ApiVersion {
	var versions []*ApiVersion
	for _, version := range a.Spec.Versions {
		versionRes := ApiVersion{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      getVersionName(a.Name, version),
				Namespace: a.Namespace,
			},
			Spec: ApiVersionSpec{
				ApiVersionSetId:   a.Status.ApiVersionSetID,
				ApiVersionScheme:  a.Spec.VersioningScheme,
				Path:              a.Spec.Path,
				APIType:           a.Spec.ApiType,
				ApiVersionSubSpec: version,
			},
		}
		versions = append(versions, &versionRes)
	}
	return versions
}

func getVersionName(apiName string, version ApiVersionSubSpec) string {
	versionSpecifier := version.Name
	if versionSpecifier == nil || *versionSpecifier == "" {
		versionSpecifier = utils.ToPointer("default")
	}
	return fmt.Sprintf("%s-%s", apiName, *versionSpecifier)
}
