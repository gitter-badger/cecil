// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package main

import (
	"fmt"
	goaclient "github.com/goadesign/goa/client"
	"github.com/spf13/cobra"
	"github.com/tleyden/cecil/goa/client"
	"github.com/tleyden/cecil/goa/tool/cli"
	"net/http"
	"os"
	"time"
)

func main() {
	// Create command line parser
	app := &cobra.Command{
		Use:   "Cecil-cli",
		Short: `CLI client for the Cecil service ()`,
	}

	// Create client struct
	httpClient := newHTTPClient()
	c := client.New(goaclient.HTTPClientDoer(httpClient))

	// Register global flags
	app.PersistentFlags().StringVarP(&c.Scheme, "scheme", "s", "", "Set the requests scheme")
	app.PersistentFlags().StringVarP(&c.Host, "host", "H", "127.0.0.1:8080", "API hostname")
	app.PersistentFlags().DurationVarP(&httpClient.Timeout, "timeout", "t", time.Duration(20)*time.Second, "Set the request timeout")
	app.PersistentFlags().BoolVar(&c.Dump, "dump", false, "Dump HTTP request and response.")

	// Register signer flags
	var key, format string
	app.PersistentFlags().StringVar(&key, "key", "", "API key used for authentication")
	app.PersistentFlags().StringVar(&format, "format", "Bearer %s", "Format used to create auth header or query from key")

	// Parse flags and setup signers
	app.ParseFlags(os.Args)
	jwtSigner := newJWTSigner(key, format)

	// Initialize API client
	c.SetJWTSigner(jwtSigner)
	c.UserAgent = "Cecil-cli/0.1"

	// Register API commands
	cli.RegisterCommands(app, c)

	// Execute!
	if err := app.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(-1)
	}
}

// newHTTPClient returns the HTTP client used by the API client to make requests to the service.
func newHTTPClient() *http.Client {
	// TBD: Change as needed (e.g. to use a different transport to control redirection policy or
	// disable cert validation or...)
	return http.DefaultClient
}

// newJWTSigner returns the request signer used for authenticating
// against the jwt security scheme.
func newJWTSigner(key, format string) goaclient.Signer {
	return &goaclient.APIKeySigner{
		SignQuery: false,
		KeyName:   "Authorization",
		KeyValue:  key,
		Format:    format,
	}

}
