package rpcserver

import (
	"context"

	"github.com/HideInBush7/go_im/pkg/redis"
	"github.com/HideInBush7/go_im/service/auth/pb"
	"github.com/HideInBush7/go_im/service/internal/model"
	"github.com/HideInBush7/go_im/service/internal/tool"
	"github.com/HideInBush7/go_im/status"
	"github.com/sirupsen/logrus"
)

func NewAuthRpcServer() *authRpcServer {
	return &authRpcServer{}
}

type authRpcServer struct {
	pb.UnimplementedAuthServer
}

// 登录
func (s *authRpcServer) Login(c context.Context, in *pb.LoginRequest) (*pb.LoginReply, error) {
	// 根据用户名查找用户
	user, err := model.GetUserByUsername(in.Username)

	if err != nil {
		// 用户不存在
		if err == model.ErrNoRows {
			return &pb.LoginReply{
				Code: status.USER_NOT_EXIST,
			}, nil
		}
		// 其他错误
		logrus.WithFields(logrus.Fields{
			"username": in.Username,
			"error":    err,
			"function": "login",
		}).Error("get user by username failed")
		return nil, err
	}

	// 密码不正确
	if user.Password != in.Password {
		return &pb.LoginReply{
			Code: status.USER_PASSWORD_WRONG,
		}, nil
	}

	// 设置token,存入redis
	token, err := SetToken(user.Uid)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"uid":      user.Uid,
			"error":    err,
			"function": "login",
		}).Error("set token failed")
		return nil, err
	}
	logrus.WithFields(logrus.Fields{
		"uid":      user.Uid,
		"username": user.Username,
	}).Debug("login success")
	return &pb.LoginReply{
		Uid:   user.Uid,
		Token: token,
		Code:  status.SUCCESS,
	}, nil
}

// 注册
func (s *authRpcServer) Register(c context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
	// 根据用户名获取token
	_, err := model.GetUserByUsername(in.Username)
	switch err {
	case nil:
		// 获取到了->用户已存在
		return &pb.RegisterReply{
			Code: status.USER_ALREADY_EXIST,
		}, nil
	case model.ErrNoRows:
		uid, err := model.InsertUser(&model.User{
			Username: in.Username,
			Password: in.Password,
		})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"username": in.Username,
				"error":    err,
				"function": "register",
			}).Error("insert user failed")
			return nil, err
		}
		token, err := SetToken(uid)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"uid":      uid,
				"username": in.Username,
				"error":    err,
				"function": "register",
			}).Error("set token failed")
			return nil, err
		}
		logrus.WithFields(logrus.Fields{
			"uid":      uid,
			"username": in.Username,
		}).Debug("login success")
		return &pb.RegisterReply{
			Uid:   uid,
			Token: token,
			Code:  status.SUCCESS,
		}, nil
	default:
		logrus.WithFields(logrus.Fields{
			"username": in.Username,
			"error":    err,
			"function": "register",
		}).Error("get user by username failed")
		return nil, err
	}
}

// 验证token
func (s *authRpcServer) Auth(c context.Context, in *pb.AuthRequest) (*pb.AuthReply, error) {
	tokenKey := tool.GetTokenKey(in.Uid, in.Token)
	ttl, err := redis.IntDo("TTL", tokenKey)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"uid":      in.Uid,
			"error":    err,
			"function": "auth",
		}).Error("ttl token failed")
		return nil, err
	}

	// key不存在
	if ttl == -2 {
		return &pb.AuthReply{
			Code: status.INVALID_TOKEN,
		}, nil
	}
	// key 存在时间 < 10分钟 重置该key过期时间为1小时
	if ttl > 0 && ttl <= 600 {
		redis.Do("SETEX", tokenKey, 3600, "")
	}
	logrus.Debug(in.Uid, " auth success")
	return &pb.AuthReply{
		Code: status.SUCCESS,
	}, nil
}

// 登出
func (s *authRpcServer) Logout(c context.Context, in *pb.LogoutRequest) (*pb.LogoutReply, error) {
	tokenKey := tool.GetTokenKey(in.Uid, in.Token)
	exists, err := redis.BoolDo("EXISTS", tokenKey)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"uid":      in.Uid,
			"error":    err,
			"function": "logout",
		}).Error("get token failed")
		return nil, err
	}
	if !exists {
		return &pb.LogoutReply{
			Code: status.INVALID_TOKEN,
		}, nil
	}
	_, err = redis.Do("DEL", tokenKey)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"uid":      in.Uid,
			"error":    err,
			"function": "logout",
		}).Error("del token failed")
		return nil, err
	}
	// TODO 断开连接层连接

	logrus.WithField("uid", in.Uid).Debug("logout success")
	return &pb.LogoutReply{
		Code: status.SUCCESS,
	}, nil
}

// 创建token,存入redis,过期时间设置为一小时
func SetToken(uid int64) (token string, err error) {
	token = tool.CreateToken(uid)
	tokenKey := tool.GetTokenKey(uid, token)
	_, err = redis.GetInstance().Do("SETEX", tokenKey, 3600, "")
	return
}
