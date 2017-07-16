@echo off
set DASH_CONFIG=./config/dashboard.json
go get
go build
dashboard.exe