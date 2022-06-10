package s3

import (
	"io"

	"github.com/Jameslikestea/d-badger/persistence"
)

var (
	_ persistence.Provider = &Service{}
	_ io.Writer            = &S3Writer{}
)

type Service struct{}

// Get implements persistence.Provider
func (*Service) Get(key string) (io.Reader, error) {
	return nil, nil
}

// Put implements persistence.Provider
func (*Service) Put(key string) (io.Writer, error) {
	return nil, nil
}

type S3Writer struct{}

// Write implements io.Writer
func (*S3Writer) Write(p []byte) (n int, err error) {
	return 0, nil
}
