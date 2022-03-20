package auth

import (
	"net"

	"github.com/HideInBush7/go_im/pkg/config"
	"github.com/HideInBush7/go_im/pkg/etcd"
	"github.com/HideInBush7/go_im/pkg/etcd/client"
	"github.com/HideInBush7/go_im/pkg/util"
	"github.com/HideInBush7/go_im/service/auth/pb"
	"github.com/HideInBush7/go_im/service/auth/rpcserver"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type authConf struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
	Rpc  string `json:"rpc"`
}

func Run() {
	// 通过配置文件加载rpc服务相关配置
	var conf = authConf{}
	config.Register("auth", &conf)

	// 默认配置
	if conf.Ip == `` {
		conf.Ip = util.GetInternalIp()
	}
	if conf.Port == `` {
		conf.Port = "8001"
	}
	if conf.Rpc == `` {
		conf.Rpc = "/auth.rpc"
	}

	// 服务运行ip:port
	addr := net.JoinHostPort(conf.Ip, conf.Port)
	lis, err := net.Listen(`tcp`, addr)
	if err != nil {
		logrus.Panic(err)
	}
	// grpc服务
	server := grpc.NewServer()
	pb.RegisterAuthServer(server, rpcserver.NewAuthRpcServer())
	go server.Serve(lis)

	// 服务注册到etcd
	etcdCli := client.GetInstance()
	r := etcd.NewRegistration(etcdCli)
	err = r.Register(conf.Rpc+addr, addr, 1)
	if err != nil {
		logrus.Panic(err)
	}
}
