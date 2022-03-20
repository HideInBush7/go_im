package main

import (
	"net/http"

	"github.com/HideInBush7/go_im/api/router"
	"github.com/HideInBush7/go_im/service/auth/rpcclient"
)

func main() {
	rpcclient.InitAuthRpcClient()
	r := router.Init()

	http.ListenAndServe(`:8000`, r)
}
