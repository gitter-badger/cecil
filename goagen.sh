#!/bin/bash

set -e
set -x

# generate API in the ./goa folder
echo "Generating API code with Goagen..."
$GOPATH/bin/goagen bootstrap -d github.com/tleyden/zerocloud/design -o ./goa
echo -e "\n"
echo    "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo    "@@@ API code generated in ./goa folder @@@"
echo    "@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@"
echo -e "\n"
