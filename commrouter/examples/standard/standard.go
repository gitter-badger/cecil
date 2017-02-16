package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/tleyden/cecil/commrouter"
)

func main() {
	cr1 := commrouter.New() // declare new CommRouter

	leaseSubject, err := cr1.AddSubject("lease") // add new subject to CommRouter
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	leaseSubject.Description = "A Lease defines the lease of an instance"

	// declare list command
	listCommand, err := leaseSubject.AddCommand("list", "show", "display") // add variations of the spelling of a command to a subject
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	listCommand.Description = "Display one or more leases"

	req := commrouter.Requirements{
		Required: true,
		Type:     commrouter.Duration,
		MinValue: time.Duration(time.Hour * 4),
	}
	err = listCommand.AddParamRequirement("this", &req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	req2 := commrouter.Requirements{
		Required:  true,
		Type:      commrouter.String,
		MustRegex: regexp.MustCompile("p([a-z]+)ch"),
	}
	err = listCommand.AddParamRequirement("fruit", &req2)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// define the controller for the list command
	listCommand.Controller(func(ctx interface{}) error {
		fmt.Println("hey! this is a new request to list!", ctx)
		return nil
	})

	// declare the terminate command
	terminateCommand, err := leaseSubject.AddCommand("terminate", "kill", "shutdown") // add variations of the spelling of a command to a subject
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	terminateCommand.Description = "Terminate a lease"
	terminateCommand.Examples = []string{"terminate lease 3"}
	// define the controller for the terminate command
	terminateCommand.Controller(func(ctx interface{}) error {
		fmt.Println("hey! this is a new request to delete!", ctx)
		return nil
	})

	err = cr1.Execute("display lease 1 selector=something param:else this:4h fruit:peach", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = cr1.Execute("terminate lease 5", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
