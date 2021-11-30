package workspace

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUnmarshalDatastore(t *testing.T) {
	a := assert.New(t)

	resp := loadExampleResp("example_resp_get_datastore.json")
	datastore := unmarshalDatastore(resp)
	a.NotEmpty(datastore)
	a.Equal("id-1", datastore.Id)
	a.Equal("datastore-1", datastore.Name)
	a.Equal("test", datastore.Description)
	a.Equal("AzureBlob", datastore.StorageType)
	a.Equal("account-1", datastore.StorageAccountName)
	a.Equal("container-1", datastore.StorageContainerName)
	a.False(datastore.IsDefault)

	// Check auth
	auth := datastore.Auth
	a.Equal("", auth.SqlUserPassword)
	a.Equal("", auth.SqlUserName)
	a.Equal("", auth.ClientId)
	a.Equal("", auth.TenantId)
	a.Equal("", auth.ClientSecret)
	a.Equal("", auth.AccountKey)
	a.Equal("AccountKey", auth.CredentialsType)

	// Check system data
	sysData := datastore.SystemData
	utcLocation, _ := time.LoadLocation("UTC")
	a.Equal(time.Date(2021, 10, 25, 10, 53, 40, 700170900, utcLocation), sysData.CreationDate)
	a.Equal("creationUser", sysData.CreationUser)
	a.Equal("Application", sysData.CreationUserType)
	a.Equal(time.Date(2021, 10, 25, 10, 53, 41, 565682100, utcLocation), sysData.LastModifiedDate)
	a.Equal("lastModifiedUser", sysData.LastModifiedUser)
	a.Equal("Application", sysData.LastModifiedUserType)
}

func TestUnmarshalDatastoreArray(t *testing.T) {
	a := assert.New(t)

	resp := loadExampleResp("example_resp_get_datastore_list.json")
	datastoreArray := unmarshalDatastoreArray(resp)
	a.NotEmpty(datastoreArray)
	a.Len(datastoreArray, 2)

	firstDatastore := datastoreArray[0]
	a.Equal("id-1", firstDatastore.Id)
	a.Equal("datastore-1", firstDatastore.Name)
	a.Equal("test", firstDatastore.Description)
	a.Equal("AzureFile", firstDatastore.StorageType)
	a.Equal("account-1", firstDatastore.StorageAccountName)
	a.Equal("container-1", firstDatastore.StorageContainerName)

	// Check auth
	auth := firstDatastore.Auth
	a.Equal("", auth.SqlUserPassword)
	a.Equal("", auth.SqlUserName)
	a.Equal("", auth.ClientId)
	a.Equal("", auth.TenantId)
	a.Equal("", auth.ClientSecret)
	a.Equal("", auth.AccountKey)
	a.Equal("AccountKey", auth.CredentialsType)

	// Check system data
	sysData := firstDatastore.SystemData
	utcLocation, _ := time.LoadLocation("UTC")
	a.Equal(time.Date(2021, 10, 7, 10, 31, 1, 714023800, utcLocation), sysData.CreationDate)
	a.Equal("redacted", sysData.CreationUser)
	a.Equal("Application", sysData.CreationUserType)
	a.Equal(time.Date(2021, 10, 7, 10, 31, 2, 649878600, utcLocation), sysData.LastModifiedDate)
	a.Equal("redacted", sysData.LastModifiedUser)
	a.Equal("Application", sysData.LastModifiedUserType)
	a.False(firstDatastore.IsDefault)
}

func TestUnmarshalDatastoreArrayEmptyResp(t *testing.T) {
	a := assert.New(t)
	resp := loadExampleResp("example_resp_get_empty_list.json")

	datastoreArray := unmarshalDatastoreArray(resp)
	a.Empty(datastoreArray)
}

func TestToWriteDatastoreSchemaEmptyDatastore(t *testing.T) {
	a := assert.New(t)
	datastore := Datastore{}
	writeSchema := toWriteDatastoreSchema(&datastore)
	a.Equal(&SchemaWrapper{WriteDatastoreSchemaProperties{Contents: WriteDatastoreSchema{}}}, writeSchema)
}

func TestToWriteDatastoreSchema(t *testing.T) {
	resp := loadExampleResp("example_resp_get_datastore.json")
	datastore := unmarshalDatastore(resp)

	writeSchema := toWriteDatastoreSchema(datastore)
	expected := &SchemaWrapper{
		WriteDatastoreSchemaProperties{
			IsDefault:   datastore.IsDefault,
			Description: datastore.Description,
			Contents: WriteDatastoreSchema{
				ContentsType:         datastore.StorageType,
				StorageAccountName:   datastore.StorageAccountName,
				StorageContainerName: datastore.StorageContainerName,
				Credentials: WriteDatastoreCredentialsSchema{
					CredentialsType: datastore.Auth.CredentialsType,
					Secrets: WriteDatastoreSecretsSchema{
						SecretsType:     datastore.Auth.CredentialsType,
						AccountKey:      datastore.Auth.AccountKey,
						ClientSecret:    datastore.Auth.ClientSecret,
						SqlUserPassword: datastore.Auth.SqlUserPassword,
					},
					ClientId:    datastore.Auth.ClientId,
					TenantId:    datastore.Auth.TenantId,
					SqlUserName: datastore.Auth.SqlUserName,
				},
			},
		},
	}
	assert.Equal(t, expected, writeSchema)
}
