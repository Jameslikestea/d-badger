package badger

import (
	"os"

	"github.com/Jameslikestea/d-badger/errors"
	"github.com/dgraph-io/badger/v3"
)

type Badger struct {
	basePath string
}

func New(bp string) *Badger {
	return &Badger{
		basePath: bp,
	}
}

func (b *Badger) Open() (*badger.DB, error) {
	dir, err := os.MkdirTemp(b.basePath, "dbadger")
	if os.IsExist(err) {
		return nil, errors.DBExist
	}

	opts := badger.DefaultOptions(dir)
	opts.Logger = nil

	return badger.Open(opts)
}

func (b *Badger) Close(bd *badger.DB) error {
	directory := bd.Opts().Dir
	return os.RemoveAll(directory)
}
