package handler

import (
	"context"
	"strconv"

	"github.com/HideInBush7/go_im/api/response"
	"github.com/HideInBush7/go_im/service/auth/pb"
	"github.com/HideInBush7/go_im/service/auth/rpcclient"
	"github.com/HideInBush7/go_im/status"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type formUser struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"username" json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var u = formUser{}
	if err := c.ShouldBindBodyWith(&u, binding.JSON); err != nil {
		response.AbortError(c, err)
		return
	}
	reply, err := rpcclient.AuthRpcClient.Register(context.Background(), &pb.RegisterRequest{
		Username: u.Username,
		Password: u.Password,
	})
	if err != nil {
		response.AbortError(c, err)
		return
	}
	if reply.Code != status.SUCCESS {
		response.Abort(c, reply.Code)
		return
	}
	response.AbortSuccess(c, map[string]interface{}{
		"uid":   reply.Uid,
		"token": reply.Token,
	})
}

func Login(c *gin.Context) {
	var u = formUser{}
	if err := c.ShouldBindBodyWith(&u, binding.JSON); err != nil {
		response.AbortError(c, err)
		return
	}
	reply, err := rpcclient.AuthRpcClient.Login(context.Background(), &pb.LoginRequest{
		Username: u.Username,
		Password: u.Password,
	})
	if err != nil {
		response.AbortError(c, err)
		return
	}
	if reply.Code != status.SUCCESS {
		response.Abort(c, reply.Code)
		return
	}
	response.AbortSuccess(c, map[string]interface{}{
		"uid":   reply.Uid,
		"token": reply.Token,
	})
}

type formLogout struct {
	Uid   string `form:"uid" json:"uid" binding:"required"`
	Token string `form:"token" json:"token" binding:"required"`
}

func Logout(c *gin.Context) {
	var l = formLogout{}
	if err := c.ShouldBindBodyWith(&l, binding.JSON); err != nil {
		response.AbortError(c, err)
		return
	}
	uid, err := strconv.ParseInt(l.Uid, 10, 64)
	if err := c.ShouldBindBodyWith(&l, binding.JSON); err != nil {
		response.AbortError(c, err)
		return
	}
	reply, err := rpcclient.AuthRpcClient.Logout(context.Background(), &pb.LogoutRequest{
		Uid:   uid,
		Token: l.Token,
	})
	if err != nil {
		response.AbortError(c, err)
		return
	}
	response.Abort(c, reply.Code)
}
