# Foobar

DBadger is a library for distributing badger KV stores across multiple concurrent services trying to access the data. It is designed to provide consistency for persistent KV data in a cost-friendly way.

## Installation

```bash
go get github.com/Jameslikestea/dbadger
```

## Usage

```go
package main

func main() {
  bs := badger.New("")
  pp := disk.New("./test")
  lm := etcd.New(c)

  lock, err := lm.Acquire("some_table")
  if err != nil {
    panic(err)
  }
  defer lm.Release(lock)

  b, _ := bs.Open()
  txn, _ := b.NewTransaction()
  txn.Set([]byte("key"), []byte("val"))
  txn.Commit()

  fh, _ := pp.Put("db1")
  _, err = b.Backup(fh)
  if err != nil {
    panic(err)
  }

}
```

## License

[Apache 2.0](https://choosealicense.com/licenses/apache-2.0/)
