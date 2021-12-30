package workspace

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestDatastorePath(t *testing.T) {
	a := assert.New(t)
	l, _ := zap.NewDevelopment()
	logger := l.Sugar()
	testCases := []struct {
		testCaseName string
		testCase     func()
	}{
		{
			testCaseName: "Test String with path containing leading /",
			testCase: func() {
				datastoreName := "datastore-1"
				datastorePath := DatastorePath{
					DatastoreName: datastoreName,
					Path:          "/foo/bar",
				}

				a.Equal("azureml://datastores/datastore-1/paths/foo/bar", datastorePath.String())
			},
		},
		{
			testCaseName: "Test String",
			testCase: func() {
				datastoreName := "datastore-1"
				datastorePath := DatastorePath{
					DatastoreName: datastoreName,
					Path:          "foo/bar",
				}

				a.Equal("azureml://datastores/datastore-1/paths/foo/bar", datastorePath.String())
			},
		},
		{
			testCaseName: "Test String empty datastore name and path",
			testCase: func() {
				datastoreName := ""
				datastorePath := DatastorePath{
					DatastoreName: datastoreName,
					Path:          "",
				}

				a.Equal("azureml://datastores//paths/", datastorePath.String())
			},
		},
		{
			testCaseName: "Test NewDatastorePath malformed path",
			testCase: func() {
				datastorePath, err := NewDatastorePath("foo/bar")
				a.Nil(datastorePath)
				a.NotNil(err)
			},
		},
		{
			testCaseName: "Test NewDatastorePath well-formed path",
			testCase: func() {
				datastoreName := "datastore1"
				path := "foo/bar/foo"
				datastorePath, err := NewDatastorePath(
					fmt.Sprintf("azureml://datastores/%s/paths/%s", datastoreName, path),
				)
				a.Nil(err)
				a.NotNil(datastorePath)
				a.Equal(datastoreName, datastorePath.DatastoreName)
				a.Equal(path, datastorePath.Path)
			},
		},
	}
	for _, t := range testCases {
		logger.Infof("Running test case %q", t.testCaseName)
		t.testCase()
	}
}
