package workspace

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net/http"
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
		responseExampleName string
		responseStatusCode  int
		error               error
		expected            *Datastore
	}{
		{
			"Get Datastore, HTTP 200",
			"example_resp_get_datastore.json",
			http.StatusOK,
			nil,
			&Datastore{
				"id-1",
				"datastore-1",
				false,
				"test",
				"account-1",
				"container-1",
				"AzureBlob",
				time.Date(2021, 10, 25, 10, 53, 40, 700170900, utcLocation),
				time.Date(2021, 10, 25, 10, 53, 41, 565682100, utcLocation),
			},
		},
	}

	for _, tc := range testCases {
		httpClient := newMockedHttpClient(tc.responseStatusCode, loadExampleResp(tc.responseExampleName), tc.error)
		workspace, _ := newClient(httpClient, &zap.SugaredLogger{})
		datastore, err := workspace.GetDatastore("foo")
		a.Equal(tc.expected, datastore, tc.description)
		a.Equal(tc.error, err, tc.description)
	}
}
