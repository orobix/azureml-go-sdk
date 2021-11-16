package workspace

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

const (
	exampleRespsDir = "assets"
)

func getExampleResp(fileName string) []byte {
	file, err := os.ReadFile(fmt.Sprintf("%s/%s", exampleRespsDir, fileName))
	if err != nil {
		log.Fatalf("Could not load example resp \"%s\"", fileName)
	}
	return file
}

func TestToDatastore(t *testing.T) {
	a := assert.New(t)

	resp := getExampleResp("example_resp_get_datastore.json")
	datastore := toDatastore(resp)
	a.NotEmpty(datastore)
	a.Equal("id-1", datastore.Id)
	a.Equal("datastore-1", datastore.Name)
	a.Equal("test", datastore.Description)
	a.Equal("AzureBlob", datastore.StorageContainerType)
	a.Equal("account-1", datastore.StorageAccountName)
	a.Equal("container-1", datastore.StorageContainerName)
	utcLocation, _ := time.LoadLocation("UTC")
	a.Equal(time.Date(2021, 10, 25, 10, 53, 40, 700170900, utcLocation), datastore.CreationDate)
	a.Equal(time.Date(2021, 10, 25, 10, 53, 41, 565682100, utcLocation), datastore.LastModifiedDate)
	a.False(datastore.IsDefault)
}

func TestToDatastoreArray(t *testing.T) {
	a := assert.New(t)

	resp := getExampleResp("example_resp_get_datastore_list.json")
	datastoreArray := toDatastoreArray(resp)
	a.NotEmpty(datastoreArray)
	a.Len(datastoreArray, 2)

	firstDatastore := datastoreArray[0]
	a.Equal("id-1", firstDatastore.Id)
	a.Equal("datastore-1", firstDatastore.Name)
	a.Equal("test", firstDatastore.Description)
	a.Equal("AzureFile", firstDatastore.StorageContainerType)
	a.Equal("account-1", firstDatastore.StorageAccountName)
	a.Equal("container-1", firstDatastore.StorageContainerName)
	utcLocation, _ := time.LoadLocation("UTC")
	a.Equal(time.Date(2021, 10, 7, 10, 31, 1, 714023800, utcLocation), firstDatastore.CreationDate)
	a.Equal(time.Date(2021, 10, 7, 10, 31, 2, 649878600, utcLocation), firstDatastore.LastModifiedDate)
	a.False(firstDatastore.IsDefault)
}

func TestToDatastoreArrayEmptyResp(t *testing.T) {
	a := assert.New(t)
	resp := getExampleResp("example_resp_get_empty_list.json")

	datastoreArray := toDatastoreArray(resp)
	a.Empty(datastoreArray)
}
