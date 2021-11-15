package aml

import (
	"fmt"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

const (
	DefaultAmlOauthScope string = "https://management.azure.com/.default"
)

type Client struct {
	workspaceHttpClient *WorkspaceHttpClient
	logger              *zap.SugaredLogger
	config              ClientConfig
}

type ClientConfig struct {
	ClientId          string
	ClientSecret      string
	TenantId          string
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
}

func NewClient(config ClientConfig, debug bool) (*Client, error) {
	credential, err := confidential.NewCredFromSecret(config.ClientSecret)
	if err != nil {
		return &Client{}, err
	}

	authority := fmt.Sprintf("https://login.microsoftonline.com/%s", config.TenantId)
	msalClient, err := confidential.New(config.ClientId, credential, confidential.WithAuthority(authority))
	if err != nil {
		return &Client{}, err
	}

	var logger *zap.Logger
	if debug == true {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}

	sugarLogger := logger.Sugar()
	client := &Client{
		workspaceHttpClient: newWorkspaceHttpClient(
			sugarLogger,
			msalClient,
			config.SubscriptionId,
			config.ResourceGroupName,
			config.WorkspaceName,
		),
		logger: sugarLogger,
	}
	return client, nil
}

func (c *Client) GetDatastores() ([]Datastore, error) {
	resp, err := c.workspaceHttpClient.doGet("datastores")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, HttpResponseError{resp.StatusCode, string(body)}
	}

	return toDatastoreArray(body), err
}

func (c *Client) GetDatastore(name string) (*Datastore, error) {
	path := fmt.Sprintf("datastores/%s", name)
	resp, err := c.workspaceHttpClient.doGet(path)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ResourceNotFoundError{"datastore", "name"}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, HttpResponseError{resp.StatusCode, string(body)}
	}

	return toDatastore(body), err
}
