package workspace

import (
	"context"
	"fmt"
	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/confidential"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type Workspace struct {
	httpClientBuilder HttpClientBuilderAPI
	logger            *zap.SugaredLogger
	datasetConverter  *DatasetConverter
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
	sugarLogger := logger.Sugar()
	return &Workspace{
		httpClientBuilder: clientBuilder,
		logger:            sugarLogger,
		datasetConverter:  &DatasetConverter{sugarLogger},
	}
}

func (w *Workspace) GetDatastores(resourceGroup, workspace string) ([]Datastore, error) {
	resp, err := w.httpClientBuilder.newClient(resourceGroup, workspace).doGet("datastores")
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

	return unmarshalDatastoreArray(body), err
}

func (w *Workspace) GetDatastore(resourceGroup, workspace, datastoreName string) (*Datastore, error) {
	path := fmt.Sprintf("datastores/%s", datastoreName)
	resp, err := w.httpClientBuilder.newClient(resourceGroup, workspace).doGet(path)
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

	return unmarshalDatastore(body), err
}

func (w *Workspace) DeleteDatastore(resourceGroup, workspace, datastoreName string) error {
	path := fmt.Sprintf("datastores/%s", datastoreName)
	resp, err := w.httpClientBuilder.newClient(resourceGroup, workspace).doDelete(path)

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

func (w *Workspace) CreateOrUpdateDatastore(resourceGroup, workspace string, datastore *Datastore) (*Datastore, error) {
	if strings.TrimSpace(datastore.Name) == "" {
		return nil, InvalidArgumentError{"the datastore name cannot be empty"}
	}

	path := fmt.Sprintf("datastores/%s", datastore.Name)
	schema := toWriteDatastoreSchema(datastore)
	resp, err := w.httpClientBuilder.newClient(resourceGroup, workspace).doPut(path, schema)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, &HttpResponseError{resp.StatusCode, string(body)}
	}

	return unmarshalDatastore(body), err
}

func (w *Workspace) CreateOrUpdateDataset(resourceGroup, workspace string, dataset *Dataset) (*Dataset, error) {
	if strings.TrimSpace(dataset.Name) == "" {
		return nil, InvalidArgumentError{"the dataset name cannot be empty"}
	}
	if len(dataset.FilePaths)+len(dataset.DirectoryPaths) == 0 {
		return nil, InvalidArgumentError{"the dataset must have at least one path"}
	}

	path := fmt.Sprintf("datasets/%s/versions/%d", dataset.Name, dataset.Version)
	schema := toWriteDatasetSchema(dataset)
	resp, err := w.httpClientBuilder.newClient(resourceGroup, workspace).doPut(path, schema)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, &HttpResponseError{resp.StatusCode, string(body)}
	}

	return w.datasetConverter.unmarshalDatasetVersion(dataset.Name, body), err
}

func (w *Workspace) GetDatasets(resourceGroup, workspace string) ([]Dataset, error) {
	datasetNames, err := w.getDatasetNames(resourceGroup, workspace)
	if err != nil {
		return nil, err
	}
	return w.retrieveLatestDatasetsVersions(resourceGroup, workspace, datasetNames)
}

func (w *Workspace) GetDatasetVersions(resourceGroup, workspace, datasetName string) ([]Dataset, error) {
	path := fmt.Sprintf("datasets/%s/versions", datasetName)
	resp, err := w.httpClientBuilder.newClient(resourceGroup, workspace).doGet(path)
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

	return w.datasetConverter.unmarshalDatasetVersionArray(datasetName, body), nil
}

// retrieveLatestDatasetsVersions For each of the dataset names provided as argument, return the respective latest version
func (w *Workspace) retrieveLatestDatasetsVersions(resourceGroup, workspaceName string, datasetNames []string) ([]Dataset, error) {
	var result []Dataset

	latestVersionChan := make(chan *Dataset, len(datasetNames))
	errChan := make(chan error, len(datasetNames))
	sem := make(chan int, NConcurrentWorkers)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	wg := sync.WaitGroup{}
	wg.Add(len(datasetNames))
	defer cancel()

	for _, datasetName := range datasetNames {
		go func(dataset string) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			case sem <- 1: // acquire lock
				d, err := w.getLatestDatasetVersion(resourceGroup, workspaceName, dataset)
				if err != nil {
					errChan <- err
					cancel()
				} else {
					latestVersionChan <- d
				}
				<-sem // release lock
			}
		}(datasetName)
	}

	wg.Wait()

	select {
	case err := <-errChan:
		return nil, err
	default:
		close(latestVersionChan)
		for d := range latestVersionChan {
			result = append(result, *d)
		}
		return result, nil
	}
}

// getLatestDatasetVersion Return the latest version of the dataset with the name provided as argument
func (w *Workspace) getLatestDatasetVersion(resourceGroup, workspace, datasetName string) (*Dataset, error) {
	w.logger.Debugf("Fetching latest version of dataset %q", datasetName)
	versions, err := w.GetDatasetVersions(resourceGroup, workspace, datasetName)
	if err != nil {
		return nil, err
	}

	var latestDataset Dataset
	for _, dataset := range versions {
		if dataset.Version > latestDataset.Version {
			latestDataset = dataset
		}
	}

	return &latestDataset, nil
}

// Return the names of the datasets of the workspace provided as argument.
func (w *Workspace) getDatasetNames(resourceGroup, workspace string) ([]string, error) {
	w.logger.Debugf("Retrieving dataset names of workspace %q in resource group %q", workspace, resourceGroup)
	resp, err := w.httpClientBuilder.newClient(resourceGroup, workspace).doGet("datasets")
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

	jsonDatastoreArray := gjson.GetBytes(body, "value").Array()
	result := make([]string, len(jsonDatastoreArray))
	for i, value := range jsonDatastoreArray {
		result[i] = value.Get("name").Str
	}
	return result, nil
}

func (w *Workspace) GetDataset(resourceGroup, workspace, name string, version int) (*Dataset, error) {
	path := fmt.Sprintf("datasets/%s/versions/%d", name, version)
	resp, err := w.httpClientBuilder.newClient(resourceGroup, workspace).doGet(path)
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

	return w.datasetConverter.unmarshalDatasetVersion(name, body), nil
}

func (w *Workspace) GetDatasetNextVersion(resourceGroup, workspace, name string) (int, error) {
	path := fmt.Sprintf("datasets/%s/versions", name)
	resp, err := w.httpClientBuilder.newClient(resourceGroup, workspace).doGet(path)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	if resp.StatusCode != http.StatusOK {
		return -1, &HttpResponseError{resp.StatusCode, string(body)}
	}

	return w.datasetConverter.unmarshalDatasetNextVersion(body), nil
}

func (w *Workspace) DeleteDataset(resourceGroup, workspace, datasetName string) error {
	path := fmt.Sprintf("datasets/%s", datasetName)
	resp, err := w.httpClientBuilder.newClient(resourceGroup, workspace).doDelete(path)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return &HttpResponseError{resp.StatusCode, string(body)}
	}

	return nil
}

func (w *Workspace) DeleteDatasetVersion(resourceGroup, workspace, datasetName string, version int) error {
	path := fmt.Sprintf("datasets/%s/versions/%d", datasetName, version)
	resp, err := w.httpClientBuilder.newClient(resourceGroup, workspace).doDelete(path)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return &HttpResponseError{resp.StatusCode, string(body)}
	}

	return nil
}
