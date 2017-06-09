// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package main

import (
	"fmt"
	"os"
	"regexp"

	. "github.com/tleyden/cecil/commrouter"
)

func main() {

	cr2 := CommRouter(func() {
		Description("this is the description of the CommRouter")
		Subject(
			"instance",
			func() {
				Description("this is the description of the subject")
				Command(
					[]string{"list", "show", "display"},
					func() {
						Description("this is the description of the command")
						Controller(func(ctx interface{}) error {
							fmt.Println("This thing works: show", ctx)
							return nil
						})

						Params(func() {
							Param("id", Int, func() {
								Required()
								MinValue(1)
							})
							Param("instance-id", String, func() {
								MinLength(1)
								MustRegex(regexp.MustCompile("i-([a-z0-9]+)"))
							})
						})

					},
				)
				Command(
					[]string{"terminate", "kill", "shutdown"},
					func() {
						Description("this is the description of the command to terminate")
						Controller(func(ctx interface{}) error {
							fmt.Println("This thing works: terminate", ctx)
							return nil
						})
					},
				)
			},
		)

		Subject(
			"pie",
			func() {
				Description("this is the description of the subject")
				Command(
					[]string{"bake", "make"},
					func() {
						Description("bake the pie")
						Controller(func(ctx interface{}) error {
							fmt.Println("This thing works: bake the pie", ctx)
							return nil
						})
					},
				)
				Command(
					[]string{"eat", "devour"},
					func() {
						Description("eat the pie")
						Controller(func(ctx interface{}) error {
							fmt.Println("This thing works: eat the pie", ctx)
							return nil
						})
					},
				)
			},
		)
	})

	err := cr2.Execute("display instance 1 selector=something id:2 instance-id:i-1g1gg1g1", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = cr2.Execute("terminate instance 1 selector=something param:else ", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = cr2.Execute("bake pie 1 selector=something param:else ", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = cr2.Execute("eat pie 1 selector=something param:else ", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
