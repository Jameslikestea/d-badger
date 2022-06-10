package main

import (
	"log"

	dbadger "github.com/Jameslikestea/d-badger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	c, err := clientv3.NewFromURLs([]string{"localhost:2379"})
	if err != nil {
		log.Fatalf("Cannot Connect To ETCD: %v", err)
	}

	s, err := session.NewSession(aws.NewConfig())
	if err != nil {
		log.Fatalf("Cannot connect to AWS: %v", err)
	}

	config := dbadger.New()
	config.WithOpts(dbadger.WithETCDLock(c), dbadger.WithS3Provider(s, "test-bucket", "badger-test"))

	d, err := config.GetDB("db1")
	if err != nil {
		log.Fatalf("Cannot get dbadger: %v", err)
	}

	b := d.Badger()
	txn := b.NewTransaction(true)
	txn.Set([]byte("hello"), []byte("world"))
	txn.Commit()

	d.Close()
}
