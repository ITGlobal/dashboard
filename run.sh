#!/bin/bash
set -e

# Build backend
echo "Building backend..."
go get
go build
chmod +x ./dashboard

# Build frontend
echo "Building frontend..."
pushd ui
npm install
npm run build:prod
popd
mkdir -p ./www
cp ./ui/dist/* ./www/

echo "Starting dashboard"
export DASH_CONFIG=./config/dashboard.json
./dashboard