package client

import (
	"time"

	"github.com/HideInBush7/go_im/pkg/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type etcdConf struct {
	EndPoints   []string `json:"endPoints"`
	DialTimeout int64    `json:"dialTimeout"`
}

var etcdClient *clientv3.Client

func init() {
	var conf = etcdConf{}
	config.Register(`etcd`, &conf)

	var err error
	etcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   conf.EndPoints,
		DialTimeout: time.Duration(conf.DialTimeout),
	})
	if err != nil {
		panic(err)
	}
}

func GetInstance() *clientv3.Client {
	return etcdClient
}
