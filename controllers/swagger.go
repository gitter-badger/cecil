// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package controllers

import "github.com/goadesign/goa"

// SwaggerController implements the swagger resource.
type SwaggerController struct {
	*goa.Controller
}

// NewSwaggerController creates a swagger controller.
func NewSwaggerController(service *goa.Service) *SwaggerController {
	return &SwaggerController{Controller: service.NewController("SwaggerController")}
}
