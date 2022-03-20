package main

import (
	"github.com/HideInBush7/go_im/pkg/log"
	"github.com/HideInBush7/go_im/service/auth/auth"
)

func main() {
	log.Init()
	auth.Run()
	select {}
}
