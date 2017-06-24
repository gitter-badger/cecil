#!/bin/bash

# Regenerate email templates
go-bindata -pkg=emailtemplates -prefix "core/email-templates/" -o emailtemplates/templates.go core/email-templates

# Regenerate go templates
go-bindata -pkg=gotemplates -prefix "core/go-templates/" -o gotemplates/templates.go core/go-templates
