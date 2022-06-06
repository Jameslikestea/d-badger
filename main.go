package main

import (
	"log"
	"os"
	"time"

	"github.com/Jameslikestea/d-badger/errors"
	"github.com/Jameslikestea/d-badger/lock"
	"github.com/Jameslikestea/d-badger/lock/etcd"
	"github.com/Jameslikestea/d-badger/persistence"
	"github.com/Jameslikestea/d-badger/persistence/disk"
	"github.com/dgraph-io/badger/v3"
	"github.com/rs/xid"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	lm lock.Manager
	pp persistence.Provider
)

func main() {
	endpoints := []string{"localhost:2379"}
	c, err := clientv3.New(clientv3.Config{
		Endpoints:        endpoints,
		DialTimeout:      5 * time.Second,
		Username:         "",
		Password:         "",
		AutoSyncInterval: time.Second,
	})
	if err != nil {
		panic("Cannot instantiate etcd connection")
	}

	lm = etcd.New(c)
	pp = disk.New("./test_lock")

	lock, err := lm.Acquire("dbadger")
	if err != nil {
		panic(err)
	}
	defer lm.Release(lock)

	dir := "C:\\Users\\james\\AppData\\Local\\Temp\\dbadger_3216119078"

	log.Printf("Opening: %s\n", dir)
	db, err := badger.Open(badger.DefaultOptions(dir))
	if err != nil {
		log.Fatal(err)
	}

	rdr, err := pp.Get("3216.bbk")
	if err != nil {
		if err != errors.NoSuchKey && err != errors.KeyNotBlob {
			panic(err)
		}
	} else {
		err = db.Load(rdr, 1)
		if err != nil {
			log.Fatal(err)
		}
	}

	id := xid.New()

	txn := db.NewTransaction(true)
	e := badger.NewEntry(id.Bytes(), []byte("hello world")).WithTTL(10 * time.Second)
	txn.SetEntry(e)
	txn.Set([]byte("grmpkg"), []byte{1})
	txn.Commit()

	var dst []byte

	txn = db.NewTransaction(false)
	i, err := txn.Get([]byte("grmpkg"))
	dst, err = i.ValueCopy(nil)

	log.Println(dst)

	wtr, err := pp.Put("3217.bbk")
	if err != nil {
		log.Println(err.Error())
	}

	_, err = db.Backup(wtr, 0)
	if err != nil {
		log.Println(err)
	}

	db.Close()
	os.RemoveAll(dir)

	time.Sleep(30 * time.Second)
}
