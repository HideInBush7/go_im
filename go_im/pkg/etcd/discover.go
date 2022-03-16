package etcd

// 服务发现

import (
	"context"
	"fmt"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

// 实现grpc 的resolver.Builder接口
func NewBuilder(c *clientv3.Client) resolver.Builder {
	return &etcdBuilder{
		client: c,
	}
}

type etcdBuilder struct {
	client *clientv3.Client
}

func (b *etcdBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	ctx, cancel := context.WithCancel(context.Background())

	r := &etcdResolver{
		c:  b.client,
		cc: cc,

		ctx:    ctx,
		cancel: cancel,
		target: target.URL.Path,
	}

	// 核心:Resolver UpdateStat更新服务地址
	go r.watch()

	return r, nil
}

func (b *etcdBuilder) Scheme() string {
	return "etcd"
}

// grpc resolver.Target 解释:
// eg: "dns://some_authority/foo.bar"
// URL.Schema: dns
// URL.Host: some_authority		=> Authority
// URP.PATH: foo.bar			=> Endpoint
// 此处约定 URL.Host为空, URL.Scheme="etcd", URL.PATH=服务名称+本地[ip:port]
// 即客户端调用时Dial etcd:///server127.0.0.1:8001
// (ip地址不应该用环回地址,不然可能会产生冲突,此处仅为示例)

// 实现grpc resolver.Resolver接口
type etcdResolver struct {
	c  *clientv3.Client
	cc resolver.ClientConn

	ctx    context.Context
	cancel context.CancelFunc
	target string // grpc resolver.Target.URL.Path
}

func (r *etcdResolver) ResolveNow(resolver.ResolveNowOptions) {}

func (r *etcdResolver) Close() {
	r.cancel()
}

func (r *etcdResolver) watch() {
	// 前缀查找etcd中target为前缀的kv对
	resp, err := r.c.Get(r.ctx, r.target, clientv3.WithPrefix())
	if err != nil {
		return
	}
	// map形式方便管理 连接断开时直接delete
	addrs := make(map[string]resolver.Address)
	for _, kv := range resp.Kvs {
		addrs[string(kv.Key)] = resolver.Address{
			Addr: string(kv.Value),
		}
	}
	// 第一次全量更新
	r.cc.UpdateState(resolver.State{
		Addresses: convertToGRPCAddresses(addrs),
	})

	wChan := r.c.Watch(r.ctx, r.target, clientv3.WithPrefix())
	for {
		select {
		case <-r.ctx.Done():
			return
		case wRes, ok := <-wChan:
			if !ok {
				fmt.Println(`resolver: watch closed`)
				return
			}
			if wRes.Err() != nil {
				fmt.Println(`resolver: watch failed `, err.Error())
				return
			}
			if len(wRes.Events) > 0 {
				for _, e := range wRes.Events {
					switch e.Type {
					// 监听到PUT
					case mvccpb.PUT:
						addrs[string(e.Kv.Key)] = resolver.Address{
							Addr: string(e.Kv.Value),
						}
					// 监听到DELETE
					case mvccpb.DELETE:
						delete(addrs, string(e.Kv.Key))
					}
				}

				// 增量更新
				r.cc.UpdateState(resolver.State{
					Addresses: convertToGRPCAddresses(addrs),
				})
			}
		}
	}
}

// 将map转换为slice
func convertToGRPCAddresses(addrs map[string]resolver.Address) []resolver.Address {
	var result []resolver.Address
	for _, v := range addrs {
		result = append(result, v)
	}
	return result
}
