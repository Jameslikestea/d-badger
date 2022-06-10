package dbadger

import (
	"github.com/Jameslikestea/d-badger/dbs/badger"
	"github.com/Jameslikestea/d-badger/lock"
	"github.com/Jameslikestea/d-badger/lock/etcd"
	"github.com/Jameslikestea/d-badger/persistence"
	"github.com/Jameslikestea/d-badger/persistence/disk"
	"github.com/Jameslikestea/d-badger/persistence/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Config struct {
	Persistence persistence.Provider
	Lock        lock.Manager
	Badger      *badger.Badger
}

type Opt func(*Config)

func New() *Config {
	return &Config{
		Badger: badger.New(""),
	}
}

func (c *Config) WithOpts(opts ...Opt) {
	for _, opt := range opts {
		opt(c)
	}
}

func WithDiskProvider(path string) Opt {
	return func(c *Config) {
		c.Persistence = disk.New(path)
	}
}

func WithS3Provider(s *session.Session, bucket, prefix string) Opt {
	return func(c *Config) {
		c.Persistence = s3.New(s, bucket, prefix)
	}
}

func WithETCDLock(client *clientv3.Client) Opt {
	return func(c *Config) {
		c.Lock = etcd.New(client)
	}
}

func WithTempDirectory(path string) Opt {
	return func(c *Config) {
		c.Badger = badger.New(path)
	}
}
