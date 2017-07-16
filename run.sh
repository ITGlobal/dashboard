#!/bin/bash
set -e
export DASH_CONFIG=./config/dashboard.json
go get
go build
chmod +x ./dashboard
./dashboard