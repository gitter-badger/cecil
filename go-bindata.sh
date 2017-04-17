#!/bin/bash

# Regenerate email assets
go-bindata -pkg=emailtemplates -prefix "core/email-templates/" -o emailtemplates/templates.go core/email-templates
