package workspace

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
	"os/exec"
	"testing"
	"time"
)

func TestNewClientEmptyConfig(t *testing.T) {
	a := assert.New(t)

	client, err := New(Config{}, false)

	a.NotNil(err)
	a.Empty(client)
}

func TestNewClientInvalidAuth(t *testing.T) {
	a := assert.New(t)

	config := Config{
		ClientId:     "invalid",
		ClientSecret: "invalid",
		TenantId:     "invalid",
	}
	client, err := New(config, false)

	a.Nil(err)
	a.NotEmpty(client)
}

func TestWorkspace_GetDatastore(t *testing.T) {
	a := assert.New(t)
	utcLocation, _ := time.LoadLocation("UTC")
	testCases := []struct {
		description         string
		datastoreName       string
		responseExampleName string
		responseStatusCode  int
		getDatastoreError   error // error returned by GetDatastore
		httpClientError     error // error returned by each call of the Http Client
		expected            *Datastore
	}{
		{
			"Get Datastore, HTTP 200",
			"foo",
			"example_resp_get_datastore.json",
			http.StatusOK,
			nil,
			nil,
			&Datastore{
				Id:                   "id-1",
				Name:                 "datastore-1",
				Type:                 "AzureBlob",
				IsDefault:            false,
				Description:          "test",
				StorageAccountName:   "account-1",
				StorageContainerName: "container-1",
				StorageType:          "AzureBlob",
				Auth: &DatastoreAuth{
					CredentialsType: "AccountKey",
				},
				SystemData: &SystemData{
					CreationDate:         time.Date(2021, 10, 25, 10, 53, 40, 700170900, utcLocation),
					CreationUserType:     "Application",
					CreationUser:         "creationUser",
					LastModifiedDate:     time.Date(2021, 10, 25, 10, 53, 41, 565682100, utcLocation),
					LastModifiedUserType: "Application",
					LastModifiedUser:     "lastModifiedUser",
				},
			},
		},
		{
			"Get Datastore, HTTP 404",
			"foo",
			"example_resp_empty.json",
			http.StatusNotFound,
			&ResourceNotFoundError{"datastore", "foo"},
			nil,
			nil,
		},
		{
			"Get Datastore, HTTP != 200",
			"foo",
			"example_resp_empty.json",
			http.StatusBadRequest,
			&HttpResponseError{http.StatusBadRequest, ""},
			nil,
			nil,
		},
		{
			"Get Datastore, HTTP Client error",
			"foo",
			"example_resp_empty.json",
			http.StatusBadRequest,
			&exec.Error{
				Name: "",
				Err:  nil,
			},
			&exec.Error{
				Name: "",
				Err:  nil,
			},
			nil,
		},
	}

	for _, tc := range testCases {
		httpClient := newMockedHttpClient(
			tc.responseStatusCode,
			loadExampleResp(tc.responseExampleName),
			tc.httpClientError,
		)
		httpClientBuilder := MockedHttpClientBuilder{httpClient: httpClient}
		workspace := newWorkspace(httpClientBuilder, &zap.Logger{})
		datastore, err := workspace.GetDatastore("", "", tc.datastoreName)
		a.Equal(tc.expected, datastore, tc.description)
		a.Equal(tc.getDatastoreError, err, tc.description)
	}
}

func TestWorkspace_GetDatastores(t *testing.T) {
	a := assert.New(t)
	utcLocation, _ := time.LoadLocation("UTC")
	testCases := []struct {
		description         string
		responseExampleName string
		responseStatusCode  int
		getDatastoreError   error // error returned by GetDatastore
		httpClientError     error // error returned by each call of the Http Client
		expected            []Datastore
	}{
		{
			"Get Datastore, HTTP 200",
			"example_resp_get_datastore_list.json",
			http.StatusOK,
			nil,
			nil,
			[]Datastore{
				{
					Id:                   "id-1",
					Name:                 "datastore-1",
					Type:                 "AzureFile",
					IsDefault:            false,
					Description:          "test",
					StorageAccountName:   "account-1",
					StorageContainerName: "container-1",
					StorageType:          "AzureFile",
					Auth: &DatastoreAuth{
						CredentialsType: "AccountKey",
					},
					SystemData: &SystemData{
						CreationDate:         time.Date(2021, 10, 7, 10, 31, 1, 714023800, utcLocation),
						CreationUserType:     "Application",
						CreationUser:         "redacted",
						LastModifiedDate:     time.Date(2021, 10, 7, 10, 31, 2, 649878600, utcLocation),
						LastModifiedUserType: "Application",
						LastModifiedUser:     "redacted",
					},
				},
				{
					Id:                   "redacted",
					Name:                 "datastore-2",
					Type:                 "AzureBlob",
					IsDefault:            true,
					Description:          "",
					StorageAccountName:   "account-2",
					StorageContainerName: "container-1",
					StorageType:          "AzureBlob",

					Auth: &DatastoreAuth{
						CredentialsType: "AccountKey",
					},
					SystemData: &SystemData{
						CreationDate:         time.Date(2021, 10, 7, 10, 31, 1, 667508600, utcLocation),
						CreationUser:         "redacted",
						CreationUserType:     "Application",
						LastModifiedDate:     time.Date(2021, 10, 7, 10, 31, 2, 879810500, utcLocation),
						LastModifiedUser:     "redacted",
						LastModifiedUserType: "Application",
					},
				},
			},
		},
		{
			"HTTP 200, empty list",
			"example_resp_get_empty_list.json",
			http.StatusOK,
			nil,
			nil,
			make([]Datastore, 0),
		},
		{
			"HTTP != 200",
			"example_resp_empty.json",
			http.StatusInternalServerError,
			&HttpResponseError{http.StatusInternalServerError, ""},
			nil,
			nil,
		},
		{
			"HTTP Client error",
			"example_resp_empty.json",
			http.StatusInternalServerError,
			&exec.Error{
				Name: "",
				Err:  nil,
			},
			&exec.Error{
				Name: "",
				Err:  nil,
			},
			nil,
		},
	}

	for _, tc := range testCases {
		httpClient := newMockedHttpClient(
			tc.responseStatusCode,
			loadExampleResp(tc.responseExampleName),
			tc.httpClientError,
		)
		httpClientBuilder := MockedHttpClientBuilder{httpClient}
		workspace := newWorkspace(httpClientBuilder, &zap.Logger{})
		datastore, err := workspace.GetDatastores("", "")
		a.Equal(tc.expected, datastore, tc.description)
		a.Equal(tc.getDatastoreError, err, tc.description)
	}
}

func TestWorkspace_DeleteDatastore(t *testing.T) {
	a := assert.New(t)
	testCases := []struct {
		description        string
		datastoreName      string
		responseStatusCode int
		httpClientError    error // error returned by each call of the Http Client
		expectedError      error
	}{
		{
			"HTTP 200 OK",
			"foo",
			http.StatusOK,
			nil,
			nil,
		},
		{
			"HTTP 404 - Datastore not found",
			"foo",
			http.StatusNotFound,
			nil,
			&ResourceNotFoundError{"datastore", "foo"},
		},
		{
			"HTTP 500 - AzureML Internal error",
			"foo",
			http.StatusInternalServerError,
			nil,
			&HttpResponseError{http.StatusInternalServerError, ""},
		},
		{
			"HTTP Client error",
			"foo",
			http.StatusOK,
			&exec.Error{"", nil},
			&exec.Error{"", nil},
		},
	}

	for _, tc := range testCases {
		httpClient := newMockedHttpClient(
			tc.responseStatusCode,
			[]byte(""),
			tc.httpClientError,
		)
		builder := MockedHttpClientBuilder{httpClient}
		workspace := newWorkspace(builder, &zap.Logger{})
		err := workspace.DeleteDatastore("", "", tc.datastoreName)
		a.Equal(tc.expectedError, err, tc.description)
	}
}
