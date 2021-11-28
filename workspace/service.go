package workspace

import (
	"fmt"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type Workspace struct {
	httpClientBuilder HttpClientBuilderAPI
	logger            *zap.SugaredLogger
}

type Config struct {
	ClientId       string
	ClientSecret   string
	TenantId       string
	SubscriptionId string
}

func New(config Config, debug bool) (*Workspace, error) {
	var logger *zap.Logger
	if debug == true {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}

	credential, err := confidential.NewCredFromSecret(config.ClientSecret)
	if err != nil {
		return &Workspace{}, err
	}

	authority := fmt.Sprintf("https://login.microsoftonline.com/%s", config.TenantId)
	msalClient, err := confidential.New(config.ClientId, credential, confidential.WithAuthority(authority))
	if err != nil {
		return &Workspace{}, err
	}

	httpClientBuilder := newHttpClientBuilder(
		logger.Sugar(),
		msalClient,
		config.SubscriptionId,
	)

	return newWorkspace(httpClientBuilder, logger), nil
}

func newWorkspace(clientBuilder HttpClientBuilderAPI, logger *zap.Logger) *Workspace {
	return &Workspace{
		httpClientBuilder: clientBuilder,
		logger:            logger.Sugar(),
	}
}

func (c *Workspace) GetDatastores(resourceGroup, workspace string) ([]Datastore, error) {
	resp, err := c.httpClientBuilder.newClient(resourceGroup, workspace).doGet("datastores")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &HttpResponseError{resp.StatusCode, string(body)}
	}

	return toDatastoreArray(body), err
}

func (c *Workspace) GetDatastore(resourceGroup, workspace, datastoreName string) (*Datastore, error) {
	path := fmt.Sprintf("datastores/%s", datastoreName)
	resp, err := c.httpClientBuilder.newClient(resourceGroup, workspace).doGet(path)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, &ResourceNotFoundError{"datastore", datastoreName}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &HttpResponseError{resp.StatusCode, string(body)}
	}

	return toDatastore(body), err
}

func (c *Workspace) DeleteDatastore(resourceGroup, workspace, datastoreName string) error {
	path := fmt.Sprintf("datastores/%s", datastoreName)
	resp, err := c.httpClientBuilder.newClient(resourceGroup, workspace).doDelete(path)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		return &ResourceNotFoundError{"datastore", datastoreName}
	}
	if resp.StatusCode != http.StatusOK {
		return &HttpResponseError{resp.StatusCode, string(body)}
	}
	return nil
}

//func (c *Workspace) CreateOrUpdateDatastore() (*Datastore, error) {
//
//}
