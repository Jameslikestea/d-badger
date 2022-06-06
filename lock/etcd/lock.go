package etcd

import (
	"context"
	"fmt"

	"github.com/Jameslikestea/d-badger/lock"
	client "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

var _ lock.Manager = &Service{}

type Service struct {
	client *client.Client
}

func New(c *client.Client) *Service {
	return &Service{
		client: c,
	}
}

// Acquire implements lock.Manager
func (s *Service) Acquire(db string) (lock.Lock, error) {
	sess, err := concurrency.NewSession(s.client)
	if err != nil {
		return nil, err
	}
	l := concurrency.NewMutex(sess, fmt.Sprintf("dbadger/locks/%s", db))
	l.Lock(context.Background())
	return l, nil
}

// Release implements lock.Manager
func (s *Service) Release(l lock.Lock) error {
	l.Unlock(context.Background())
	return nil
}
