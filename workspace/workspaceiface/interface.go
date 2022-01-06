package workspaceiface

import "github.com/Telemaco019/azureml-go-sdk/workspace"

type WorkspaceAPI interface {
	// GetDatastores Return the list of datastore of the AML Workspace provided as argument.
	GetDatastores(resourceGroup, workspace string) ([]workspace.Datastore, error)

	// GetDatastore Return the datastore with the name provided as argument.
	GetDatastore(resourceGroup, workspace, datastoreName string) (*workspace.Datastore, error)

	// DeleteDatastore Delete the datastore with the name provided as argument
	DeleteDatastore(resourceGroup, workspace, datastoreName string) error

	// CreateOrUpdateDatastore Create or update the datastore with the data provided as argument
	CreateOrUpdateDatastore(resourceGroup, workspace string, datastore *workspace.Datastore) (*workspace.Datastore, error)

	// GetDatasets Return the list of datasets of the AML Workspace. For each dataset, only its latest version is returned.
	GetDatasets(resourceGroup, workspace string) ([]workspace.Dataset, error)

	// GetDataset Return the dataset with the name and version provided as argument
	GetDataset(resourceGroup, workspace, name string, version int) (*workspace.Dataset, error)

	// GetDatasetNextVersion Return the next version of the dataset with the name provided as argument
	GetDatasetNextVersion(resourceGroup, workspace, name string) (int, error)

	// GetDatasetVersions Return all the versions of the dataset with the name provided as argument
	GetDatasetVersions(resourceGroup, workspace, datasetName string) ([]workspace.Dataset, error)

	// CreateOrUpdateDataset Create or update the dataset with the data provided as argument
	CreateOrUpdateDataset(resourceGroup, workspace string, dataset *workspace.Dataset) (*workspace.Dataset, error)

	// DeleteDataset Delete the dataset (all its versions) with the name provided as argument
	DeleteDataset(resourceGroup, workspace, datasetName string) error

	// DeleteDatasetVersion Delete the version provided as argument of the dataset with the specified name
	DeleteDatasetVersion(resourceGroup, workspace, datasetName string, version int) error
}
