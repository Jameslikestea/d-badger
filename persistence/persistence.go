package persistence

import "io"

//go:generate go run github.com/golang/mock/mockgen -destination=./mocks/provider.go . Provider

type Provider interface {
	Get(key string) (io.Reader, error)
	Put(key string) (io.Writer, error)
}
