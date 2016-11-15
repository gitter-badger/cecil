#!/bin/bash

set -e
set -x

# generate API in the ./goa folder
echo "Generating API code with Goagen..."
# $GOPATH/bin/goagen bootstrap -d github.com/tleyden/cecil/design -o ./goa

$GOPATH/bin/goagen main -d github.com/tleyden/cecil/design -o ./goa-controllers
# this command will generate the contollers, but the import path for app
# will be wrong; for each of controller files in ./goa-controllers
# change "github.com/tleyden/cecil/goa-controllers/app" to
# "github.com/tleyden/cecil/goa/app"

$GOPATH/bin/goagen app -d github.com/tleyden/cecil/design -o ./goa
$GOPATH/bin/goagen client -d github.com/tleyden/cecil/design -o ./goa
$GOPATH/bin/goagen swagger -d github.com/tleyden/cecil/design -o ./goa
# $GOPATH/bin/goagen js -d github.com/tleyden/cecil/design -o ./goa
# $GOPATH/bin/goagen schema -d github.com/tleyden/cecil/design -o ./goa

echo    "API code generated in ./goa folder"
echo    "API controllers generated in ./goa-controllers folder"
