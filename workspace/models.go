package workspace

import (
	"fmt"
	"strings"
	"time"
)

const datastorePathPrefix = "azureml://datastores/"

type SystemData struct {
	CreationDate     time.Time
	CreationUser     string
	CreationUserType string

	LastModifiedDate     time.Time
	LastModifiedUser     string
	LastModifiedUserType string
}

type DatastoreAuth struct {
	CredentialsType string
	ClientId        string
	TenantId        string
	ClientSecret    string
	AccountKey      string
	SqlUserName     string
	SqlUserPassword string
}

type Datastore struct {
	Id          string
	Name        string
	IsDefault   bool
	Description string

	StorageType          string
	StorageAccountName   string
	StorageContainerName string

	SystemData *SystemData
	Auth       *DatastoreAuth
}

type Dataset struct {
	Id             string
	Name           string
	Description    string
	DatastoreId    string
	Version        int
	FilePaths      []DatasetPath
	DirectoryPaths []DatasetPath
	SystemData     *SystemData
}

type DatasetPath interface {
	fmt.Stringer
}

type DatastorePath struct {
	DatastoreName string
	Path          string
}

func NewDatastorePath(path string) (*DatastorePath, error) {
	datastoreNameWithPath := strings.TrimPrefix(path, datastorePathPrefix)
	parts := strings.Split(datastoreNameWithPath, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf(
			"malformed path, datastore path sould be in the format %s/<datastore-name>/paths/<path>",
			datastorePathPrefix,
		)
	}
	datastoreName := parts[0]
	return &DatastorePath{
		DatastoreName: datastoreName,
		Path:          strings.Join(parts[2:], "/"),
	}, nil
}

func (d DatastorePath) String() string {
	var cleanedPath string
	if len(d.Path) > 0 && d.Path[0:1] == "/" {
		cleanedPath = d.Path[1:]
	} else {
		cleanedPath = d.Path
	}
	return fmt.Sprintf("azureml://datastores/%s/paths/%s", d.DatastoreName, cleanedPath)
}
