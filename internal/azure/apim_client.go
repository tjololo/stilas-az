package azure

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	apim "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v2"
	"net/http"
)

// APIMClient is a client for interacting with the Azure API Management service
type APIMClient struct {
	// ApimClientConfig is the configuration for the APIM client
	ApimClientConfig  ApimClientConfig
	apimClientFactory *apim.ClientFactory
}

// ApimClientConfig is the configuration for the APIMClient
type ApimClientConfig struct {
	ClientOptions   *azidentity.DefaultAzureCredentialOptions
	FactoryOptions  *arm.ClientOptions
	SubscriptionId  string
	ResourceGroup   string
	ApimServiceName string
}

// NewAPIMClient creates a new APIMClient
func NewAPIMClient(config ApimClientConfig) (*APIMClient, error) {
	credential, err := azidentity.NewDefaultAzureCredential(config.ClientOptions)
	if err != nil {
		return nil, err
	}
	clientFactory, err := apim.NewClientFactory(config.SubscriptionId, credential, config.FactoryOptions)
	if err != nil {
		return nil, err
	}
	return &APIMClient{
		ApimClientConfig:  config,
		apimClientFactory: clientFactory,
	}, nil
}

func (c *APIMClient) GetApiVersionSet(ctx context.Context, apiVersionSetName string, options *apim.APIVersionSetClientGetOptions) (apim.APIVersionSetClientGetResponse, error) {
	client := c.apimClientFactory.NewAPIVersionSetClient()
	return client.Get(ctx, c.ApimClientConfig.ResourceGroup, c.ApimClientConfig.ApimServiceName, apiVersionSetName, options)
}

func (c *APIMClient) CreateUpdateApiVersionSet(ctx context.Context, apiVersionSetName string, parameters apim.APIVersionSetContract, options *apim.APIVersionSetClientCreateOrUpdateOptions) (apim.APIVersionSetClientCreateOrUpdateResponse, error) {
	client := c.apimClientFactory.NewAPIVersionSetClient()
	return client.CreateOrUpdate(ctx, c.ApimClientConfig.ResourceGroup, c.ApimClientConfig.ApimServiceName, apiVersionSetName, parameters, options)
}

func (c *APIMClient) GetApi(ctx context.Context, apiId string, options *apim.APIClientGetOptions) (apim.APIClientGetResponse, error) {
	client := c.apimClientFactory.NewAPIClient()
	return client.Get(ctx, c.ApimClientConfig.ResourceGroup, c.ApimClientConfig.ApimServiceName, apiId, options)
}

func (c *APIMClient) CreateUpdateApi(ctx context.Context, apiId string, parameters apim.APICreateOrUpdateParameter, options *apim.APIClientBeginCreateOrUpdateOptions) (*runtime.Poller[apim.APIClientCreateOrUpdateResponse], error) {
	client := c.apimClientFactory.NewAPIClient()
	return client.BeginCreateOrUpdate(ctx, c.ApimClientConfig.ResourceGroup, c.ApimClientConfig.ApimServiceName, apiId, parameters, options)
}

func IsNotFoundError(err error) bool {
	var responseError *azcore.ResponseError
	if errors.As(err, &responseError) {
		return responseError.StatusCode == http.StatusNotFound
	}
	return false
}

func IgnoreNotFound(err error) error {
	if IsNotFoundError(err) {
		return nil
	}
	return err
}
