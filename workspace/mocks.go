package workspace

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
)

type MockedHttpClientBuilder struct {
	httpClient HttpClientAPI
}

func (m MockedHttpClientBuilder) newClient(_, _ string) HttpClientAPI {
	return m.httpClient
}

type MockedHttpClient struct {
	mock.Mock
}

// expected args:
// - position 0: the response status code (int)
// - position 1: the response body (string)
// - position 2: the returned error (error)
func (t MockedHttpClient) doGet(path string) (*http.Response, error) {
	args := t.Called(path)
	mockedResponse := &http.Response{
		StatusCode: args.Int(0),
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(args.String(1)))),
	}
	return mockedResponse, args.Error(2)
}

// expected args:
// - position 0: the response status code (int)
// - position 1: the response body (string)
// - position 2: the returned error (error)
func (t MockedHttpClient) doGetWithContext(ctx context.Context, path string) (*http.Response, error) {
	args := t.Called(ctx, path)
	mockedResponse := &http.Response{
		StatusCode: args.Int(0),
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(args.String(1)))),
	}
	return mockedResponse, args.Error(2)
}

// expected args:
// - position 0: the response status code (int)
// - position 1: the response body (string)
// - position 2: the returned error (error)
func (t MockedHttpClient) doDelete(path string) (*http.Response, error) {
	args := t.Called(path)
	mockedResponse := &http.Response{
		StatusCode: args.Int(0),
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(args.String(1)))),
	}
	return mockedResponse, args.Error(2)
}

// expected args:
// - position 0: the response status code (int)
// - position 1: the response body (string)
// - position 2: the returned error (error)
func (t MockedHttpClient) doPut(path string, requestBody interface{}) (*http.Response, error) {
	args := t.Called(path, requestBody)
	mockedResponse := &http.Response{
		StatusCode: args.Int(0),
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(args.String(1)))),
	}
	return mockedResponse, args.Error(2)
}
