package main

import (
	"log"

	dbadger "github.com/Jameslikestea/d-badger"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	c, err := clientv3.NewFromURLs([]string{"localhost:2379"})
	if err != nil {
		log.Fatalf("Cannot Connect To ETCD: %v", err)
	}

	config := dbadger.New()
	config.WithOpts(dbadger.WithETCDLock(c), dbadger.WithDiskProvider("./test_lock"))

	d, err := config.GetDB("db1")
	if err != nil {
		log.Fatalf("Cannot get dbadger: %v", err)
	}
	defer d.Close()

	b := d.Badger()
	txn := b.NewTransaction(true)
	txn.Set([]byte("hello"), []byte("world"))
	txn.Commit()
}
