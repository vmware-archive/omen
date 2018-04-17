package sessions

import (
	"time"
)

//go:generate counterfeiter . opsmanClient
type opsmanClient interface {
	Delete(endpoint string, timeout time.Duration) error
}

func NewSessionManager(client opsmanClient) SessionManager {
	return SessionManager{
		opsmanClient: client,
	}
}

type SessionManager struct {

	opsmanClient opsmanClient

}

func (manager SessionManager) ClearAll() error {

	return manager.opsmanClient.Delete("/api/v0/sessions", 5*time.Minute)

}