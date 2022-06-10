package dbadger

import (
	"github.com/Jameslikestea/d-badger/errors"
	"github.com/Jameslikestea/d-badger/lock"
	"github.com/dgraph-io/badger/v3"
)

type Connection struct {
	l    lock.Lock
	b    *badger.DB
	name string
	conf *Config
}

func (c *Config) GetDB(name string) (*Connection, error) {
	l, err := c.Lock.Acquire(name)
	if err != nil {
		return nil, err
	}

	db, err := c.Badger.Open()
	if err != nil {
		c.Lock.Release(l)
		return nil, err
	}

	bk, err := c.Persistence.Get(name)
	if err != nil {
		if err != errors.NoSuchKey && err != errors.KeyNotBlob {
			c.Lock.Release(l)
			return nil, err
		}
	}
	if bk != nil {
		err = db.Load(bk, 512)
		if err != nil {
			c.Lock.Release(l)
			return nil, err
		}
	}

	return &Connection{
		l:    l,
		b:    db,
		name: name,
		conf: c,
	}, nil
}

func (c *Connection) Badger() *badger.DB {
	return c.b
}

func (c *Connection) Close() error {
	defer c.conf.Lock.Release(c.l)
	wrtr, err := c.conf.Persistence.Put(c.name)
	if err != nil {
		return err
	}
	_, err = c.b.Backup(wrtr, 0)
	if err != nil {
		return err
	}

	return nil
}
