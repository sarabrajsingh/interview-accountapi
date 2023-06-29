package mocks

import "net/http"

// TO-DO implement a mock-client that emulated http.Client.Do for integration tets
type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

var (
	GetDoFunc func(req *http.Request) (*http.Response, error)
)
