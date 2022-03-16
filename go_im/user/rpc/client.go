package rpc

import (
	"context"
	"fmt"

	"github.com/HideInBush7/go_im/pb"
	"github.com/HideInBush7/go_im/pkg/etcd"
	"github.com/HideInBush7/go_im/pkg/etcd/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserRpcClient() *userRpcClient {
	etcdClient := client.GetInstance()
	builder := etcd.NewBuilder(etcdClient)

	conn, err := grpc.Dial(`etcd:///user.rpc`,
		grpc.WithResolvers(builder),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
	)
	if err != nil {
		panic(err)
	}
	return &userRpcClient{
		conn: pb.NewUserClient(conn),
	}
}

type userRpcClient struct {
	conn pb.UserClient
}

func (u *userRpcClient) Register(username, password string) (code int, token string, msg string) {
	reply, err := u.conn.Register(context.Background(), &pb.RegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return int(pb.Code_ERROR), ``, err.Error()
	}
	code = int(reply.Code)
	token = reply.Token
	msg = reply.Code.String()
	return
}
