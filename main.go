package main

import (
	"flag"
	"fmt"

	"github.com/inconshreveable/log15"
	"github.com/tleyden/zerocloud/core"
)

func main() {

	flag.BoolVar(&core.DropAllTables, "drop-all-tables", false, "If passed, drops all tables")
	flag.Parse()

	if core.DropAllTables {
		fmt.Println("You are about to drop all tables from DB; are you sure? [N/y]")
		isSure := core.AskForConfirmation()
		if isSure {
			fmt.Println("Tables WILL BE dropped.")
		} else {
			fmt.Println("Tables will NOT be dropped.")
			core.DropAllTables = false
		}
	}

	service := core.NewService()
	// service will be passed to Goa controllers

	service.SetupAndRun()
	defer service.Stop(true)

	if err := service.RunHTTPServer(); err != nil {
		log15.Error("service.RunHTTPServer()", "error", err)
	}
}
