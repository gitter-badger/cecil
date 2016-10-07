package main

import (
	"github.com/tleyden/zerocloud/core"
)

func main() {
	service := core.Initialize()
	defer service.Stop()
	service.RunHTTPServer()
}
