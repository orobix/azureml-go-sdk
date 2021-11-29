package workspace

import (
	"bytes"
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

func (m *MockedHttpClient) doDelete(path string) (*http.Response, error) {
	return &m.mockedResponse, m.err
}

func (m *MockedHttpClient) doPut(path string, body interface{}) (*http.Response, error) {
	return &m.mockedResponse, m.err
}
