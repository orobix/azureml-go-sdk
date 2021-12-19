package workspace

import (
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
	}
	for _, t := range testCases {
		logger.Infof("Running test case %q", t.testCaseName)
		t.testCase()
	}
}
