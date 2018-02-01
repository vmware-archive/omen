package fakes

import "time"

type FakeOMClient struct {
	GetFunc  func(string) ([]byte, error)
	PostFunc func(string, string) ([]byte, error)
}

func (f FakeOMClient) Get(endpoint string, timeout time.Duration) ([]byte, error) {
	return f.GetFunc(endpoint)
}

func (f FakeOMClient) Post(endpoint, data string, timeout time.Duration) ([]byte, error) {
	return f.PostFunc(endpoint, data)
}
