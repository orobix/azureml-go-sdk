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

func newMockedHttpClient(responseStatusCode int, responseBody []byte, err error) *mockedHttpClient {
	bodyReader := ioutil.NopCloser(bytes.NewReader(responseBody))
	mockedResponse := http.Response{
		StatusCode: responseStatusCode,
		Body:       bodyReader,
	}
	return &mockedHttpClient{mockedResponse, err}
}

func (m *mockedHttpClient) doGet(path string) (*http.Response, error) {
	return &m.mockedResponse, m.err
}
