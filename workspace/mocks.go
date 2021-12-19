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

func newMockedHttpClient(responseStatusCode int, responseBody []byte, err error) *MockedHttpClient {
	bodyReader := ioutil.NopCloser(bytes.NewReader(responseBody))
	mockedResponse := http.Response{
		StatusCode: responseStatusCode,
		Body:       bodyReader,
	}
	return &MockedHttpClient{mockedResponse, err}
}

type MockedHttpClient struct {
	mockedResponse http.Response
	err            error
}

func (m *MockedHttpClient) doGet(path string) (*http.Response, error) {
	return &m.mockedResponse, m.err
}

func (m *MockedHttpClient) doGetWithContext(ctx context.Context, path string) (*http.Response, error) {
	return &m.mockedResponse, m.err
}

func (m *MockedHttpClient) doDelete(path string) (*http.Response, error) {
	return &m.mockedResponse, m.err
}

func (m *MockedHttpClient) doPut(path string, body interface{}) (*http.Response, error) {
	return &m.mockedResponse, m.err
}

type TestifyMockedHttpClient struct {
	mock.Mock
}

// expected args:
// - position 0: the response status code (int)
// - position 1: the response body (string)
// - position 2: the returned error (error)
func (t TestifyMockedHttpClient) doGet(path string) (*http.Response, error) {
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
func (t TestifyMockedHttpClient) doGetWithContext(ctx context.Context, path string) (*http.Response, error) {
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
func (t TestifyMockedHttpClient) doDelete(path string) (*http.Response, error) {
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
func (t TestifyMockedHttpClient) doPut(path string, requestBody interface{}) (*http.Response, error) {
	args := t.Called(path, requestBody)
	mockedResponse := &http.Response{
		StatusCode: args.Int(0),
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(args.String(1)))),
	}
	return mockedResponse, args.Error(2)
}
