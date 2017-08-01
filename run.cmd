@echo off
rem Build backend
echo Building backend...
go get
if %ERRORLEVEL% NEQ 0 (
    echo Command "go get" exited with %ERRORLEVEL%
    exit 1
)
go build
if %ERRORLEVEL% NEQ 0 (
    echo Command "go build" exited with %ERRORLEVEL%
    exit 1
)

rem Build frontend
echo Building frontend...
cd ui

call npm install
if %ERRORLEVEL% NEQ 0 (
    echo Command "npm install" exited with %ERRORLEVEL%
    exit 1
)

call npm run build:prod
if %ERRORLEVEL% NEQ 0 (
    echo Command "npm run build:prod" exited with %ERRORLEVEL%
    exit 1
)

cd ..
if not exist ./www mkdir ./www
copy /Y ui\dist\* www
if %ERRORLEVEL% NEQ 0 (
    echo Command "copy /Y ui\dist\* www" exited with %ERRORLEVEL%
    exit 1
)

rem Run application
echo Starting dashboard
set DASH_CONFIG=./config/dashboard.json
dashboard.exe