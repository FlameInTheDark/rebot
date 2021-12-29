package etcd

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewETCDClient(endpoint string) (*clientv3.Client, error) {
	return clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: 5 * time.Second,
	})
}
