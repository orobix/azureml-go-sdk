package workspace

import "time"

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

	SystemData SystemData
	Auth       DatastoreAuth
}
