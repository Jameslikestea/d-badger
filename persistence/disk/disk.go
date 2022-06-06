package disk

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/Jameslikestea/d-badger/errors"
	"github.com/Jameslikestea/d-badger/persistence"
)

var _ persistence.Provider = &Service{}

type Service struct {
	basePath string
}

func New(p string) *Service {
	i, err := os.Stat(p)
	if err == os.ErrNotExist || err == os.ErrExist {
		err = os.MkdirAll(p, 0755)
		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	} else if !i.IsDir() {
		panic("Not a directory")
	}

	return &Service{
		basePath: p,
	}
}

// Get implements persistence.Provider
func (s *Service) Get(key string) (io.Reader, error) {
	keyFile := filepath.Join(s.basePath, key)

	// Check some preconditions with the information provided
	i, err := os.Stat(keyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.NoSuchKey
		}
		return nil, errors.KeyNotBlob
	}
	if i.IsDir() {
		return nil, errors.KeyNotBlob
	}

	readStream, err := os.Open(keyFile)
	if err != nil {
		return nil, errors.KeyError
	}

	return readStream, nil
}

// Put implements persistence.Provider
func (s *Service) Put(key string) (io.Writer, error) {
	keyFile := filepath.Join(s.basePath, key)

	// Check some preconditions with the information provided
	i, err := os.Stat(keyFile)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Println(err)
			return nil, errors.DestError
		}
	} else {
		if i.IsDir() {
			return nil, errors.DestIsDirectory
		}
	}

	writeStream, err := os.OpenFile(keyFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	return writeStream, nil
}
