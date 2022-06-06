package lock

//go:generate go run github.com/golang/mock/mockgen -destination ./mocks/lock.go -package=mocks . Lock,Manager

import "context"

// Lock
// This is an abstraction layer above any vendor specific lock (in this case etcd) all vendor specific
// implementation should be written in a concrete struct that implements this interface.
type Lock interface {
	Lock(ctx context.Context) error
	Unlock(ctx context.Context) error
}

// Manager
// This interface provides a few functions that should be implemented in any locking service (the vendor
// specific implementation of a service)
type Manager interface {
	Acquire(db string) (Lock, error)
	Release(Lock) error
}
