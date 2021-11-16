package workspace

import "time"

type Datastore struct {
	Id                   string
	Name                 string
	IsDefault            bool
	Description          string
	StorageAccountName   string
	StorageContainerName string
	StorageContainerType string
	CreationDate         time.Time
	LastModifiedDate     time.Time
}
