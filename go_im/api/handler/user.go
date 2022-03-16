package handler

import (
	"fmt"

	"github.com/HideInBush7/go_im/user/rpc"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	cli := rpc.NewUserRpcClient()
	code, token, msg := cli.Register(`hello`, `world`)
	if code != 0 {
		fmt.Println(token, msg)
	}
}
