package main

import (
	"github.com/tleyden/zerocloud/core"
)

func main() {
	service := core.NewService()
	service.Setup()
	defer service.Stop()
	service.RunHTTPServer()
}
