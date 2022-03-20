package etcd

// 服务注册

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func NewRegistration(c *clientv3.Client) *etcdRegistration {
	return &etcdRegistration{
		client: c,
	}
}

type etcdRegistration struct {
	client  *clientv3.Client
	leaseId clientv3.LeaseID // 服务注册的租约ID
}

// 服务注册
// target: 服务名称 对应grpc resolver.Target.URL.Path
// addr: 服务真实地址
// ttl: 服务续约时间(s) 建议设短一点,否则服务断开后etcd一段时间仍有该租约
func (r *etcdRegistration) Register(target string, addr string, ttl int64) error {
	// 创建租约
	lease, err := r.client.Grant(context.Background(), ttl)
	if err != nil {
		return err
	}
	r.leaseId = lease.ID

	// 使用该租约注册
	_, err = r.client.Put(context.Background(), target+addr, addr, clientv3.WithLease(r.leaseId))
	if err != nil {
		return err
	}

	// 自动续租
	keepAlive, err := r.client.KeepAlive(context.Background(), r.leaseId)
	if err != nil {
		return err
	}

	// 开启消费
	go func() {
		for {
			// close(keepAlive)会<-nil
			if res := <-keepAlive; res == nil {
				fmt.Println(`keepAlive closed...`)
				return
			}
		}
	}()
	return nil
}

// 服务主动注销
func (r *etcdRegistration) Unregister() error {
	_, err := r.client.Revoke(context.Background(), r.leaseId)
	return err
}
