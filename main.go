package main

import (
	"github.com/inconshreveable/log15"
	"github.com/tleyden/zerocloud/core"
)

func main() {
	service := core.NewService()
	// service will be passed to Goa controllers

	service.SetupAndRun()
	defer service.Stop(true)

	if err := service.RunHTTPServer(); err != nil {
		log15.Error("service.RunHTTPServer()", "error", err)
	}
}
