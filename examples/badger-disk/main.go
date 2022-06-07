package main

import (
	"log"
	"time"

	"github.com/Jameslikestea/d-badger/dbs/badger"
	"github.com/Jameslikestea/d-badger/errors"
	"github.com/Jameslikestea/d-badger/lock"
	"github.com/Jameslikestea/d-badger/lock/etcd"
	"github.com/Jameslikestea/d-badger/persistence"
	"github.com/Jameslikestea/d-badger/persistence/disk"
	badge "github.com/dgraph-io/badger/v3"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	lm lock.Manager
	pp persistence.Provider
	db *badger.Badger
)

func main() {
	c, err := clientv3.NewFromURLs([]string{"localhost:2379"})
	if err != nil {
		log.Fatalf("Cannot Connect To ETCD: %v", err)
	}

	lm = etcd.New(c)
	pp = disk.New("./test_lock")
	// Use an empty string to use the default temp directory of the OS
	db = badger.New("")

	// Acquire a lock
	l, err := lm.Acquire("db1")
	if err != nil {
		log.Fatalf("Cannot Acquire Lock: %v", err)
	}
	defer lm.Release(l)

	// Open a new empty database
	b, err := db.Open()
	if err != nil {
		log.Fatalf("Cannot Open Database: %v", err)
	}

	// Get the previous backup (This is the bit that we need to lock for)
	// in this case we want to allow for an empty/fresh database to be
	// loaded
	bbk, err := pp.Get("db1")
	if err != nil && err != errors.NoSuchKey && err != errors.KeyNotBlob {
		log.Fatalf("Cannot Get Previous Backup: %v", err)
	}

	if bbk != nil {
		// Choose a random number
		b.Load(bbk, 1024)
	}

	txn := b.NewTransaction(true)
	ent := badge.NewEntry([]byte("some-key"), []byte("some-val")).WithTTL(10 * time.Minute)
	txn.SetEntry(ent)
	err = txn.Commit()
	if err != nil {
		log.Fatalf("Could not set entry: %v", err)
	}

	writer, err := pp.Put("db1")
	if err != nil {
		log.Fatalf("Cannot update database: %v", err)
	}
	_, err = b.Backup(writer, 0)
	if err != nil {
		log.Fatalf("Cannot Backup Data: %v", err)
	}
}
