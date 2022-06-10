# DBadger

DBadger is a library for distributing badger KV stores across multiple concurrent services trying to access the data. It is designed to provide consistency for persistent KV data in a cost-friendly way.

## Installation

```bash
go get github.com/Jameslikestea/dbadger
```

## Usage

See [Examples](./examples/) for a showcase of some of the configurations this library supports.

```go
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
```

## License

[Apache 2.0](https://choosealicense.com/licenses/apache-2.0/)
