package main

import (
	"gServ/core/config"
	"gServ/core/gameserv"
	"gServ/core/httpserv"
	"gServ/core/log"
	"gServ/core/repository"
	"gServ/core/tcpserv"
	"gServ/core/validate"
)

func init() {
	if err := config.Init("gserv"); err != nil {
		panic(err)
	}

	if err := log.Init(config.GetConfig().Server.Log); err != nil {
		panic(err)
	}

	if err := validate.Init(); err != nil {
		panic(err)
	}

	if err := repository.Init(); err != nil {
		panic(err)
	}

	if err := gameserv.Init(); err != nil {
		panic(err)
	}

	if err := httpserv.Init(); err != nil {
		panic(err)
	}

	if err := tcpserv.Init(); err != nil {
		panic(err)
	}
}
