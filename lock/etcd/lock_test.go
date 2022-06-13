package etcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/server/v3/embed"
	"testing"
)

func TestLock(t *testing.T) {
	e := startEtcd(t)
	defer e.Server.Stop()

	c, err := clientv3.NewFromURLs(e.Server.Cfg.ClientURLs.StringSlice())
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	lm := New(c)

	lock, err := lm.Acquire("test-lock")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	err = lm.Release(lock)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func startEtcd(t *testing.T) *embed.Etcd {
	cfg := embed.NewConfig()

	cfg.Dir = "__test__.etcd"
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	<-e.Server.ReadyNotify()

	return e
}
