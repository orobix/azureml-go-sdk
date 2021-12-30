package workspace

import (
	"fmt"
	"github.com/tidwall/gjson"
	"regexp"
)

func unmarshalDatastoreArray(json []byte) []Datastore {
	jsonDatastoreArray := gjson.GetBytes(json, "value").Array()
	datastoreSlice := make([]Datastore, gjson.GetBytes(json, "value.#").Int())
	for i, jsonDatastore := range jsonDatastoreArray {
		datastore := unmarshalDatastore([]byte(jsonDatastore.Raw))
		datastoreSlice[i] = *datastore
	}
	return datastoreSlice
}

func unmarshalDatastore(json []byte) *Datastore {
	auth := DatastoreAuth{
		CredentialsType: gjson.GetBytes(json, "properties.contents.credentials.credentialsType").Str,
		TenantId:        gjson.GetBytes(json, "properties.contents.credentials.tenantId").Str,
		ClientId:        gjson.GetBytes(json, "properties.contents.credentials.clientId").Str,
		ClientSecret:    gjson.GetBytes(json, "properties.contents.credentials.secret.clientSecret").Str,
		AccountKey:      gjson.GetBytes(json, "properties.contents.credentials.secret.accountKey").Str,
		SqlUserName:     gjson.GetBytes(json, "properties.contents.credentials.secret.userId").Str,
		SqlUserPassword: gjson.GetBytes(json, "properties.contents.credentials.secret.password").Str,
	}
	return &Datastore{
		Id:                   gjson.GetBytes(json, "id").Str,
		Name:                 gjson.GetBytes(json, "name").Str,
		Description:          gjson.GetBytes(json, "properties.description").Str,
		IsDefault:            gjson.GetBytes(json, "properties.isDefault").Bool(),
		StorageAccountName:   gjson.GetBytes(json, "properties.contents.accountName").Str,
		StorageContainerName: gjson.GetBytes(json, "properties.contents.containerName").Str,
		StorageType:          gjson.GetBytes(json, "properties.contents.contentsType").Str,

		SystemData: unmarshalSystemData(json),
		Auth:       &auth,
	}
}

func unmarshalDatasetVersionArray(datasetName string, json []byte) []Dataset {
	jsonDatasetArray := gjson.GetBytes(json, "value").Array()
	datasetSlice := make([]Dataset, gjson.GetBytes(json, "value.#").Int())
	for i, jsonDataset := range jsonDatasetArray {
		dataset := unmarshalDatasetVersion(datasetName, []byte(jsonDataset.Raw))
		datasetSlice[i] = *dataset
	}
	return datasetSlice
}

func unmarshalDatasetVersion(datasetName string, json []byte) *Dataset {
	return &Dataset{
		Id:             gjson.GetBytes(json, "id").Str,
		Name:           datasetName,
		Description:    gjson.GetBytes(json, "properties.description").Str,
		DatastoreId:    gjson.GetBytes(json, "properties.datastoreId").Str,
		Version:        int(gjson.GetBytes(json, "name").Int()),
		FilePaths:      unmarshalDatasetPaths(gjson.GetBytes(json, "properties.paths"), "file"),
		DirectoryPaths: unmarshalDatasetPaths(gjson.GetBytes(json, "properties.paths"), "folder"),
		SystemData:     unmarshalSystemData(json),
	}
}

func unmarshalDatasetPaths(jsonDatasetPaths gjson.Result, pathType string) []DatasetPath {
	result := make([]DatasetPath, 0)
	jsonDatasetPaths.ForEach(func(key, value gjson.Result) bool {
		path := value.Get(pathType)
		if path.Exists() == false {
			return false // TODO: log error
		}
		if path.Type != gjson.Null {
			isDatastorePath, _ := regexp.MatchString(fmt.Sprintf("%s.*", datastorePathPrefix), path.Str)
			if isDatastorePath == true {
				datastorePath, err := NewDatastorePath(path.Str)
				if err != nil {
					// TODO: log error
				} else {
					result = append(result, datastorePath)
				}
			}
		}
		return true
	})
	return result
}

func unmarshalSystemData(json []byte) *SystemData {
	return &SystemData{
		CreationDate:         gjson.GetBytes(json, "systemData.createdAt").Time(),
		CreationUserType:     gjson.GetBytes(json, "systemData.createdByType").Str,
		CreationUser:         gjson.GetBytes(json, "systemData.createdBy").Str,
		LastModifiedDate:     gjson.GetBytes(json, "systemData.lastModifiedAt").Time(),
		LastModifiedUserType: gjson.GetBytes(json, "systemData.lastModifiedByType").Str,
		LastModifiedUser:     gjson.GetBytes(json, "systemData.lastModifiedBy").Str,
	}
}

func toWriteDatastoreSchema(datastore *Datastore) *SchemaWrapper {
	var secrets *WriteDatastoreSecretsSchema
	var credentials *WriteDatastoreCredentialsSchema

	if datastore.Auth != nil {
		secrets = &WriteDatastoreSecretsSchema{
			SecretsType:     datastore.Auth.CredentialsType,
			AccountKey:      datastore.Auth.AccountKey,
			ClientSecret:    datastore.Auth.ClientSecret,
			SqlUserPassword: datastore.Auth.SqlUserPassword,
		}
		credentials = &WriteDatastoreCredentialsSchema{
			CredentialsType: datastore.Auth.CredentialsType,
			Secrets:         secrets,
			ClientId:        datastore.Auth.ClientId,
			TenantId:        datastore.Auth.TenantId,
			SqlUserName:     datastore.Auth.SqlUserName,
		}
	}

	return &SchemaWrapper{
		Properties: WriteDatastoreSchemaProperties{
			IsDefault:   datastore.IsDefault,
			Description: datastore.Description,
			Contents: WriteDatastoreSchema{
				ContentsType:         datastore.StorageType,
				StorageAccountName:   datastore.StorageAccountName,
				StorageContainerName: datastore.StorageContainerName,
				Credentials:          credentials,
			},
		},
	}
}

func toWriteDatasetSchema(dataset *Dataset) *SchemaWrapper {
	pathSchemas := make([]DatasetPathsSchema, len(dataset.FilePaths)+len(dataset.DirectoryPaths))
	for i, filePath := range dataset.FilePaths {
		pathSchemas[i] = DatasetPathsSchema{FilePath: filePath.String()}
	}
	for i, directoryPath := range dataset.DirectoryPaths {
		pathSchemas[i] = DatasetPathsSchema{DirectoryPath: directoryPath.String()}
	}

	return &SchemaWrapper{
		Properties: WriteDatasetSchema{
			Description: dataset.Description,
			Paths:       pathSchemas,
		},
	}
}
