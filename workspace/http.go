package workspace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"go.uber.org/zap"
	"io"
	"net/http"
)

const (
	amlApiVersion          = "2021-10-01"
	amlWorkspaceApiBaseUrl = "https://management.azure.com/subscriptions/%s/resourceGroups/%s/providers/Microsoft.MachineLearningServices/workspaces/%s"
)

type HttpClientBuilderAPI interface {
	newClient(resourceGroupName, workspaceName string) HttpClientAPI
}

func newHttpClientBuilder(
	logger *zap.SugaredLogger,
	msalClient confidential.Client,
	subscriptionId string) HttpClientBuilderAPI {
	return &HttpClientBuilder{
		logger:         logger,
		msalClient:     msalClient,
		subscriptionId: subscriptionId,
		httpClient:     &http.Client{},
	}
}

type HttpClientBuilder struct {
	logger         *zap.SugaredLogger
	msalClient     confidential.Client
	subscriptionId string
	httpClient     *http.Client
}

func (b *HttpClientBuilder) newClient(resourceGroupName, workspaceName string) HttpClientAPI {
	return &HttpClient{
		logger:            b.logger,
		msalClient:        b.msalClient,
		subscriptionId:    b.subscriptionId,
		resourceGroupName: resourceGroupName,
		workspaceName:     workspaceName,
		httpClient:        b.httpClient,
	}
}

type HttpClientAPI interface {
	doGet(path string) (*http.Response, error)

	doGetWithContext(ctx context.Context, path string) (*http.Response, error)

	doDelete(path string) (*http.Response, error)

	doPut(path string, requestBody interface{}) (*http.Response, error)
}

type HttpClient struct {
	logger            *zap.SugaredLogger
	msalClient        confidential.Client
	subscriptionId    string
	resourceGroupName string
	workspaceName     string
	httpClient        *http.Client
}

func (c HttpClient) getJwt() (string, error) {
	scopes := []string{DefaultAmlOauthScope}
	c.logger.Debug("Using cached JWT silently...")
	authResult, err := c.msalClient.AcquireTokenSilent(context.Background(), scopes)
	if err != nil {
		c.logger.Debug("Could not acquire JWT silently, now acquiring it with Client Credential flow...")
		authResult, err = c.msalClient.AcquireTokenByCredential(context.Background(), scopes)
		c.logger.Debug("JWT acquired")
	}
	return authResult.AccessToken, err
}

func (c *HttpClient) getWorkspaceApiBaseUrl() string {
	return fmt.Sprintf(amlWorkspaceApiBaseUrl, c.subscriptionId, c.resourceGroupName, c.workspaceName)
}

func (c *HttpClient) prepareRequest(req *http.Request) error {
	jwt, err := c.getJwt()
	if err != nil {
		return err
	}

	// Add required headers
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))

	// Add required query params
	q := req.URL.Query()
	q.Add("api-version", amlApiVersion)
	req.URL.RawQuery = q.Encode()
	return nil
}

func (c *HttpClient) newRequest(method string, url string, requestBody []byte) (*http.Request, error) {
	var requestBodyReader io.Reader
	if requestBody == nil {
		requestBodyReader = nil
	} else {
		requestBodyReader = bytes.NewBuffer(requestBody)
	}

	req, err := http.NewRequest(method, url, requestBodyReader)
	if err != nil {
		return req, err
	}

	err = c.prepareRequest(req)
	return req, err
}

func (c *HttpClient) newRequestWithContext(ctx context.Context, method string, url string, requestBody []byte) (*http.Request, error) {
	var requestBodyReader io.Reader
	if requestBody == nil {
		requestBodyReader = nil
	} else {
		requestBodyReader = bytes.NewBuffer(requestBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, requestBodyReader)
	if err != nil {
		return req, err
	}

	err = c.prepareRequest(req)
	return req, err
}

func (c *HttpClient) doGet(path string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.getWorkspaceApiBaseUrl(), path)
	request, err := c.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	c.logger.Infof("GET > %s", request.URL)
	return c.httpClient.Do(request)
}

func (c HttpClient) doGetWithContext(ctx context.Context, path string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.getWorkspaceApiBaseUrl(), path)
	request, err := c.newRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	c.logger.Infof("GET > %s", request.URL)
	return c.httpClient.Do(request)
}

func (c *HttpClient) doDelete(path string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.getWorkspaceApiBaseUrl(), path)
	request, err := c.newRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	c.logger.Infof("DELETE > %s", request.URL)
	return c.httpClient.Do(request)
}

func (c *HttpClient) doPut(path string, requestBody interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.getWorkspaceApiBaseUrl(), path)

	b, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	request, err := c.newRequest("PUT", url, b)
	request.Header.Add("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	c.logger.Infof("PUT > %s", request.URL)
	return c.httpClient.Do(request)
}
