package main

import (
	"fmt"
	"net"

	"github.com/HideInBush7/go_im/pb"
	"github.com/HideInBush7/go_im/pkg/etcd"
	"github.com/HideInBush7/go_im/pkg/etcd/client"
	"github.com/HideInBush7/go_im/pkg/util"
	"github.com/HideInBush7/go_im/user/rpc"
	"google.golang.org/grpc"
)

func main() {
	addr := net.JoinHostPort(util.GetInternalIp(), "8001")
	lis, _ := net.Listen(`tcp`, addr)
	s := grpc.NewServer()
	pb.RegisterUserServer(s, rpc.NewUserRpcServer())
	go s.Serve(lis)

	cli := client.GetInstance()
	reg := etcd.NewRegistration(cli)
	err := reg.Register("/user.rpc"+addr, addr, 10)
	fmt.Println(err)
	select {}
}
