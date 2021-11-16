package workspace

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type mockedHttpClient struct {
	mockedResponse http.Response
	err            error
}

func NewMockedHttpClient(responseStatusCode int, responseBody string, err error) *mockedHttpClient {
	bodyReader := ioutil.NopCloser(bytes.NewReader([]byte(responseBody)))
	mockedResponse := http.Response{
		StatusCode: responseStatusCode,
		Body:       bodyReader,
	}
	return &mockedHttpClient{mockedResponse, err}
}

func (m *mockedHttpClient) doGet(path string) (*http.Response, error) {
	return &m.mockedResponse, m.err
}
