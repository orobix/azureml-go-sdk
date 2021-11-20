package workspace

import (
	"fmt"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type Workspace struct {
	httpClient HttpClientAPI
	logger     *zap.SugaredLogger
	config     Config
}

type Config struct {
	ClientId          string
	ClientSecret      string
	TenantId          string
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
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

	httpClient := newHttpClient(
		logger.Sugar(),
		msalClient,
		config.SubscriptionId,
		config.ResourceGroupName,
		config.WorkspaceName,
	)

	return newClient(httpClient, logger.Sugar())
}

func newClient(httpClient HttpClientAPI, logger *zap.SugaredLogger) (*Workspace, error) {
	client := &Workspace{
		httpClient: httpClient,
		logger:     logger,
	}
	return client, nil
}

func (c *Workspace) GetDatastores() ([]Datastore, error) {
	resp, err := c.httpClient.doGet("datastores")
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

func (c *Workspace) GetDatastore(name string) (*Datastore, error) {
	path := fmt.Sprintf("datastores/%s", name)
	resp, err := c.httpClient.doGet(path)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, &ResourceNotFoundError{"datastore", name}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, &HttpResponseError{resp.StatusCode, string(body)}
	}

	return toDatastore(body), err
}

func (c *Workspace) DeleteDatastore(name string) error {
	path := fmt.Sprintf("datastores/%s", name)
	resp, err := c.httpClient.doDelete(path)

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		return &ResourceNotFoundError{"datastore", name}
	}
	if resp.StatusCode != http.StatusOK {
		return &HttpResponseError{resp.StatusCode, string(body)}
	}
	return nil
}
