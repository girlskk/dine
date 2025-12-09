package zxh

import (
	"sync"

	"github.com/go-resty/resty/v2"
)

type Manager struct {
	pool sync.Pool
}

func NewManager(baseURL string) *Manager {
	rc := resty.New().
		SetBaseURL(baseURL).
		SetHeader("Content-Type", "application/json")

	return &Manager{
		pool: sync.Pool{
			New: func() any { return &Client{rc: rc} },
		},
	}
}

func (m *Manager) GetClient(mac *Credentials) (*Client, func()) {
	client := m.pool.Get().(*Client)
	client.mac = mac
	return client, func() { client.mac = nil; m.pool.Put(client) }
}
