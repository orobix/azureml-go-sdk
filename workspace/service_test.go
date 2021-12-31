package workspace

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		mockedHttpClient := new(MockedHttpClient)
		mockedHttpClient.On("doGet", mock.Anything).Return(
			tc.responseStatusCode,
			string(loadExampleResp(tc.responseExampleName)),
			tc.httpClientError,
		)
		httpClientBuilder := MockedHttpClientBuilder{httpClient: mockedHttpClient}
		logger, _ := zap.NewDevelopment()
		workspace := newWorkspace(httpClientBuilder, logger)
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
		mockedHttpClient := new(MockedHttpClient)
		mockedHttpClient.On("doGet", mock.Anything).Return(
			tc.responseStatusCode,
			string(loadExampleResp(tc.responseExampleName)),
			tc.httpClientError,
		)
		httpClientBuilder := MockedHttpClientBuilder{mockedHttpClient}
		logger, _ := zap.NewDevelopment()
		workspace := newWorkspace(httpClientBuilder, logger)

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
		mockedHttpClient := new(MockedHttpClient)
		mockedHttpClient.On("doDelete", mock.Anything).Return(
			tc.responseStatusCode,
			"",
			tc.httpClientError,
		)
		builder := MockedHttpClientBuilder{mockedHttpClient}
		logger, _ := zap.NewDevelopment()
		workspace := newWorkspace(builder, logger)
		err := workspace.DeleteDatastore("", "", tc.datastoreName)
		a.Equal(tc.expectedError, err, tc.description)
	}
}

func TestWorkspace_CreateOrUpdateDatastore(t *testing.T) {
	a := assert.New(t)
	testCases := []struct {
		description         string
		inputDatastore      *Datastore
		responseStatusCode  int
		responseExampleName string
		httpClientError     error // error returned by each call of the Http Client
		expectedError       error
	}{
		{
			"Invalid input: datastore without name",
			&Datastore{},
			http.StatusOK,
			"example_resp_empty.json",
			nil,
			InvalidArgumentError{"the datastore name cannot be empty"},
		},
		{
			"HTTP 201 - Created",
			&Datastore{Name: "foo"},
			http.StatusCreated,
			"example_resp_get_datastore.json",
			nil,
			nil,
		},
		{
			"HTTP 500 - AzureML Internal error",
			&Datastore{Name: "foo"},
			http.StatusInternalServerError,
			"example_resp_empty.json",
			nil,
			&HttpResponseError{http.StatusInternalServerError, ""},
		},
		{
			"HTTP Client error",
			&Datastore{Name: "foo"},
			http.StatusOK,
			"example_resp_empty.json",
			&exec.Error{"", nil},
			&exec.Error{"", nil},
		},
	}

	for _, tc := range testCases {
		mockedHttpClient := new(MockedHttpClient)
		mockedHttpClient.On("doPut", mock.Anything, mock.Anything).Return(
			tc.responseStatusCode,
			string(loadExampleResp(tc.responseExampleName)),
			tc.httpClientError,
		)
		builder := MockedHttpClientBuilder{mockedHttpClient}
		logger, _ := zap.NewDevelopment()
		workspace := newWorkspace(builder, logger)
		ds, err := workspace.CreateOrUpdateDatastore("", "", tc.inputDatastore)

		if err == nil {
			a.Equal(unmarshalDatastore(loadExampleResp(tc.responseExampleName)), ds)
		} else {
			a.Nil(ds)
			a.Equal(tc.expectedError, err, tc.description)
		}
	}
}

func TestWorkspace_RetrieveLatestDatasetsVersions(t *testing.T) {
	a := assert.New(t)
	l, _ := zap.NewDevelopment()
	logger := l.Sugar()
	testCases := []struct {
		testCaseName        string
		testCaseDescription string
		testCase            func()
	}{
		{
			testCaseName: "Test empty input dataset names array",
			testCase: func() {
				mockedHttpClient := new(MockedHttpClient)
				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				latestVersions, err := ws.retrieveLatestDatasetsVersions("", "", []string{})
				a.Nil(err)
				a.Empty(latestVersions)
			},
		},
		{
			testCaseName: "Test resp in error in fetching latest dataset version",
			testCase: func() {
				mockedResponseBody := "error"
				mockedResponseStatusCode := http.StatusInternalServerError
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doGet", mock.Anything).Return(mockedResponseStatusCode, mockedResponseBody, nil)
				mockedDatasetList := getMockedDatasetNames(NConcurrentWorkers * 2)

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				latestVersions, err := ws.retrieveLatestDatasetsVersions("", "", mockedDatasetList)
				a.Nil(latestVersions)
				a.Equal(&HttpResponseError{mockedResponseStatusCode, mockedResponseBody}, err)
			},
		},
		{
			testCaseName: "Test http client returns error in fetching latest dataset version",
			testCase: func() {
				mockedResponseBody := ""
				mockedResponseStatusCode := http.StatusInternalServerError
				mockedError := fmt.Errorf("error")
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doGet", mock.Anything).Return(mockedResponseStatusCode, mockedResponseBody, mockedError)
				mockedDatasetList := getMockedDatasetNames(NConcurrentWorkers * 2)

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				latestVersions, err := ws.retrieveLatestDatasetsVersions("", "", mockedDatasetList)
				a.Nil(latestVersions)
				a.Equal(mockedError, err)
			},
		},
		{
			testCaseName: "Test retrieve latest datasets versions success",
			testCase: func() {
				mockedResponseBody := string(loadExampleResp("example_resp_get_dataset_versions.json"))
				mockedResponseStatusCode := http.StatusOK
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doGet", mock.Anything).Return(mockedResponseStatusCode, mockedResponseBody, nil)
				mockedDatasetList := getMockedDatasetNames(NConcurrentWorkers * 2)

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				latestVersions, err := ws.retrieveLatestDatasetsVersions("", "", mockedDatasetList)
				a.NotEmpty(latestVersions)
				a.Nil(err)
			},
		},
	}
	for _, testCase := range testCases {
		logger.Infof("Running test case %q", testCase.testCaseName)
		testCase.testCase()
	}
}

func TestWorkspace_GetLatestDatasetVersion(t *testing.T) {
	a := assert.New(t)
	l, _ := zap.NewDevelopment()
	logger := l.Sugar()
	testCases := []struct {
		testCaseName        string
		testCaseDescription string
		testCase            func()
	}{
		{
			testCaseName: "Test get latest version empty resp",
			testCase: func() {
				mockedResponseBody := ""
				mockedResponseStatusCode := http.StatusOK
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doGet", mock.Anything).Return(mockedResponseStatusCode, mockedResponseBody, nil)

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				latestVersion, err := ws.getLatestDatasetVersion("", "", "")
				a.Nil(err)
				a.Empty(latestVersion)
			},
		},
		{
			testCaseName: "Test get latest version http response is in error",
			testCase: func() {
				mockedResponseBody := "error"
				mockedResponseStatusCode := http.StatusInternalServerError
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doGet", mock.Anything).Return(mockedResponseStatusCode, mockedResponseBody, nil)

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				latestVersion, err := ws.getLatestDatasetVersion("", "", "")
				a.Empty(latestVersion)
				a.Equal(&HttpResponseError{mockedResponseStatusCode, mockedResponseBody}, err)
			},
		},
		{
			testCaseName: "Test get latest version success",
			testCase: func() {
				mockedResponseBody := string(loadExampleResp("example_resp_get_dataset_versions.json"))
				mockedResponseStatusCode := http.StatusOK
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doGet", mock.Anything).Return(mockedResponseStatusCode, mockedResponseBody, nil)

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				latestVersion, err := ws.getLatestDatasetVersion("", "", "")
				a.Nil(err)
				a.Equal(4, latestVersion.Version)
			},
		},
	}

	for _, testCase := range testCases {
		logger.Infof("Running test case %q", testCase.testCaseName)
		testCase.testCase()
	}
}

func TestWorkspace_getDatasetNames(t *testing.T) {
	a := assert.New(t)
	l, _ := zap.NewDevelopment()
	logger := l.Sugar()
	testCases := []struct {
		testCaseName        string
		testCaseDescription string
		testCase            func()
	}{
		{
			testCaseName: "Test get dataset names http response is in error",
			testCase: func() {
				mockedResponseBody := "error"
				mockedResponseStatusCode := http.StatusInternalServerError
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doGet", mock.Anything).Return(mockedResponseStatusCode, mockedResponseBody, nil)

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				names, err := ws.getDatasetNames("", "")
				a.Empty(names)
				a.Equal(&HttpResponseError{mockedResponseStatusCode, mockedResponseBody}, err)
			},
		},
		{
			testCaseName: "Test get dataset names http client returns error",
			testCase: func() {
				mockedResponseBody := "error"
				mockedError := fmt.Errorf("error")
				mockedResponseStatusCode := http.StatusInternalServerError
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doGet", mock.Anything).Return(mockedResponseStatusCode, mockedResponseBody, mockedError)

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				names, err := ws.getDatasetNames("", "")
				a.Empty(names)
				a.Equal(mockedError, err)
			},
		},
		{
			testCaseName: "Test get dataset names success",
			testCase: func() {
				mockedResponseBody := string(loadExampleResp("example_resp_get_datasets.json"))
				mockedResponseStatusCode := http.StatusOK
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doGet", mock.Anything).Return(mockedResponseStatusCode, mockedResponseBody, nil)

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				names, err := ws.getDatasetNames("", "")
				a.Nil(err)
				a.Len(names, 3)
			},
		},
	}

	for _, testCase := range testCases {
		logger.Infof("Running test case %q", testCase.testCaseName)
		testCase.testCase()
	}
}

func TestWorkspace_CreateOrUpdateDataset(t *testing.T) {
	a := assert.New(t)
	l, _ := zap.NewDevelopment()
	logger := l.Sugar()
	testCases := []struct {
		testCaseName        string
		testCaseDescription string
		testCase            func()
	}{
		{
			testCaseName: "Test create or update dataset with empty name",
			testCase: func() {
				mockedHttpClient := new(MockedHttpClient)
				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				d := &Dataset{Name: "          "}
				updatedDataset, err := ws.CreateOrUpdateDataset("", "", d)
				a.Nil(updatedDataset)
				a.NotEmpty(err)
			},
		},
		{
			testCaseName: "Test create or update dataset http response 404 not found",
			testCase: func() {
				mockedResponseBody := "error"
				mockedResponseStatusCode := http.StatusNotFound
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doPut", mock.Anything, mock.Anything).Return(mockedResponseStatusCode, mockedResponseBody, nil)

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				latestVersion, err := ws.CreateOrUpdateDataset("", "", &Dataset{Name: "foo"})
				a.Empty(latestVersion)
				a.Equal(&HttpResponseError{mockedResponseStatusCode, mockedResponseBody}, err)
			},
		},
		{
			testCaseName: "Test create or update dataset http client returns error",
			testCase: func() {
				mockedResponseBody := "error"
				mockedError := fmt.Errorf("error")
				mockedResponseStatusCode := http.StatusInternalServerError
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doPut", mock.Anything, mock.Anything).Return(mockedResponseStatusCode, mockedResponseBody, mockedError)

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				names, err := ws.CreateOrUpdateDataset("", "", &Dataset{Name: "foo"})
				a.Empty(names)
				a.Equal(mockedError, err)
			},
		},
		{
			testCaseName: "Test create or update dataset success",
			testCase: func() {
				mockedResponseBody := string(loadExampleResp("example_resp_create_or_update_dataset.json"))
				mockedResponseStatusCode := http.StatusOK
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doPut", mock.Anything, mock.Anything).Return(mockedResponseStatusCode, mockedResponseBody, nil)

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				updatedDataset, err := ws.CreateOrUpdateDataset("", "", &Dataset{Name: "foo"})
				a.Nil(err)
				a.Equal("redacted", updatedDataset.Id)
				a.Equal(1, len(updatedDataset.DirectoryPaths))
				a.Equal(0, len(updatedDataset.FilePaths))
			},
		},
	}

	for _, testCase := range testCases {
		logger.Infof("Running test case %q", testCase.testCaseName)
		testCase.testCase()
	}
}

func TestWorkspace_GetDataset(t *testing.T) {
	a := assert.New(t)
	l, _ := zap.NewDevelopment()
	logger := l.Sugar()
	testCases := []struct {
		testCaseName        string
		testCaseDescription string
		testCase            func()
	}{
		{
			testCaseName: "Test get dataset not found resp",
			testCase: func() {
				mockedResponseBody := "not found"
				mockedResponseStatusCode := http.StatusNotFound
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doGet", mock.Anything).Return(mockedResponseStatusCode, mockedResponseBody, nil)

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				dataset, err := ws.GetDataset("", "", "", 1)
				a.Nil(dataset)
				a.Equal(&HttpResponseError{mockedResponseStatusCode, mockedResponseBody}, err)
			},
		},
		{
			testCaseName: "Test get dataset http response is in error",
			testCase: func() {
				mockedResponseBody := "error"
				mockedResponseStatusCode := http.StatusInternalServerError
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doGet", mock.Anything).Return(mockedResponseStatusCode, mockedResponseBody, nil)

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				dataset, err := ws.GetDataset("rg", "ws", "dataset", 1)
				a.Empty(dataset)
				a.Equal(&HttpResponseError{mockedResponseStatusCode, mockedResponseBody}, err)
			},
		},
		{
			testCaseName: "Test get dataset http client is in error",
			testCase: func() {
				clientErrorMsg := "error"
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doGet", mock.Anything).Return(1, "", fmt.Errorf(clientErrorMsg))

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				dataset, err := ws.GetDataset("rg", "ws", "dataset", 1)
				a.Nil(dataset)
				a.Equal(clientErrorMsg, err.Error())
			},
		},
		{
			testCaseName: "Test get dataset success",
			testCase: func() {
				mockedResponseBody := string(loadExampleResp("example_resp_get_dataset.json"))
				mockedResponseStatusCode := http.StatusOK
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doGet", mock.Anything).Return(mockedResponseStatusCode, mockedResponseBody, nil)

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				dataset, err := ws.GetDataset("rg", "ws", "dataset", 1)
				a.Nil(err)
				a.Equal("<id>", dataset.Id)
				a.Equal(1, len(dataset.FilePaths))
				a.NotEmpty(dataset.SystemData)
			},
		},
	}

	for _, testCase := range testCases {
		logger.Infof("Running test case %q", testCase.testCaseName)
		testCase.testCase()
	}
}
func TestWorkspace_GetDatasetNextVersion(t *testing.T) {
	a := assert.New(t)
	l, _ := zap.NewDevelopment()
	logger := l.Sugar()
	testCases := []struct {
		testCaseName        string
		testCaseDescription string
		testCase            func()
	}{
		{
			testCaseName: "Test get dataset next version not found resp",
			testCase: func() {
				mockedResponseBody := "not found"
				mockedResponseStatusCode := http.StatusNotFound
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doGet", mock.Anything).Return(mockedResponseStatusCode, mockedResponseBody, nil)

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				nextVersion, err := ws.GetDatasetNextVersion("rg", "ws", "dataset")
				a.Equal(-1, nextVersion)
				a.Equal(&HttpResponseError{mockedResponseStatusCode, mockedResponseBody}, err)
			},
		},
		{
			testCaseName: "Test get dataset next version http response is in error",
			testCase: func() {
				mockedResponseBody := "error"
				mockedResponseStatusCode := http.StatusInternalServerError
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doGet", mock.Anything).Return(mockedResponseStatusCode, mockedResponseBody, nil)

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				nextVersion, err := ws.GetDatasetNextVersion("rg", "ws", "dataset")
				a.Equal(-1, nextVersion)
				a.Equal(&HttpResponseError{mockedResponseStatusCode, mockedResponseBody}, err)
			},
		},
		{
			testCaseName: "Test get dataset next version http client is in error",
			testCase: func() {
				clientErrorMsg := "error"
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doGet", mock.Anything).Return(1, "", fmt.Errorf(clientErrorMsg))

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				nextVesion, err := ws.GetDatasetNextVersion("rg", "ws", "dataset")
				a.Equal(-1, nextVesion)
				a.Equal(clientErrorMsg, err.Error())
			},
		},
		{
			testCaseName: "Test get dataset next version success",
			testCase: func() {
				mockedResponseBody := string(loadExampleResp("example_resp_get_dataset_next_version.json"))
				mockedResponseStatusCode := http.StatusOK
				mockedHttpClient := new(MockedHttpClient)
				mockedHttpClient.On("doGet", mock.Anything).Return(mockedResponseStatusCode, mockedResponseBody, nil)

				builder := MockedHttpClientBuilder{mockedHttpClient}
				ws := newWorkspace(builder, l)
				nextVersion, err := ws.GetDatasetNextVersion("rg", "ws", "dataset")
				a.Nil(err)
				a.Equal(8, nextVersion)
			},
		},
	}

	for _, testCase := range testCases {
		logger.Infof("Running test case %q", testCase.testCaseName)
		testCase.testCase()
	}
}

func getMockedDatasetNames(n int) []string {
	result := make([]string, n)
	for i := 0; i < n; i++ {
		result[i] = fmt.Sprintf("dataset-%d", i)
	}
	return result
}
