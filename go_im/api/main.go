package main

import (
	"net/http"

	"github.com/HideInBush7/go_im/api/router"
)

func main() {
	r := router.Init()

	http.ListenAndServe(`:8000`, r)
}
