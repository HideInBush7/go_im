package rpc

import (
	"context"

	"github.com/HideInBush7/go_im/pb"
	"github.com/sirupsen/logrus"
)

func NewUserRpcServer() pb.UserServer {
	return &userRpcServer{}
}

type userRpcServer struct {
	pb.UnimplementedUserServer
}

// 登录
func (u *userRpcServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginReply, error) {
	panic(``)
}

// 注册
func (u *userRpcServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
	logrus.Info(in.Username)
	logrus.Info(in.Password)
	return &pb.RegisterReply{
		Token: "hello",
		Code:  pb.Code_SUCCESS,
	}, nil
}

// 根据uid获取用户
func (u *userRpcServer) GetUserByUid(ctx context.Context, in *pb.GetUserByUidRequest) (*pb.GetUserByUidReply, error) {
	panic(``)
}

// 验证token
func (u *userRpcServer) Auth(ctx context.Context, in *pb.AuthRequest) (*pb.AuthReply, error) {
	panic(``)
}

// 登出
func (u *userRpcServer) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutReply, error) {
	panic(``)
}
