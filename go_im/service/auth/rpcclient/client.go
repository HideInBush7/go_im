package rpcclient

import (
	"fmt"
	"sync"

	"github.com/HideInBush7/go_im/pkg/etcd"
	"github.com/HideInBush7/go_im/pkg/etcd/client"
	"github.com/HideInBush7/go_im/service/auth/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

var once sync.Once

var AuthRpcClient pb.AuthClient

func InitAuthRpcClient() {
	once.Do(func() {
		// 获取etcd client
		cli := client.GetInstance()
		// 服务发现 WithResolvers
		builder := etcd.NewBuilder(cli)
		conn, err := grpc.Dial("etcd:///auth.rpc",
			grpc.WithResolvers(builder),
			// 负载均衡
			grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy":"%s"}`, roundrobin.Name)),
			grpc.WithInsecure(),
		)
		if err != nil {
			panic(err)
		}
		AuthRpcClient = pb.NewAuthClient(conn)
	})
}
