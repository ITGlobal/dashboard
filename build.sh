#!/bin/bash

BUILD_VERSION=""

make_version() {
    BUILD_VERSION=$(git describe --tags)
    echo Building IT Global Dashboard v$BUILD_VERSION
}

build() {
    GOARCH=$1
    GOOS=$2

    mkdir -p ./out/$2

    go build -o ./out/$2/$3 -ldflags "-X main.Version=v$BUILD_VERSION -X main.BuildConfiguration=[$2]"
    zip ./out/dashboard-$2.zip ./out/$2/$3

    echo Compiled Dashboard v$BUILD_VERSION for $2/$1 - see ./out/dashboard-$2.zip   
}

build_win32() {
    build "386" "windows" "win32" "dashboard.exe"
}

build_win64() {
    build "amd64" "windows" "win64" "dashboard.exe"
}

build_linux32() {
    build "386" "linux" "linux32" "dashboard"
}

build_linux64() {
    build "amd64" "linux" "linux64" "dashboard"
}

build_linuxarm() {
    build "arm" "linux" "linuxarm" "dashboard"
}


if [[ "$1" -eq "" ]]; then
	make_version
    build_win32
    build_win64
    build_linux32
    build_linux64
    build_linuxarm
    exit
fi


if [[ "$1" -eq "win32" ]]; then
	make_version
    build_win32
    exit
fi

if [[ "$1" -eq "win64" ]]; then
	make_version
    build_win64
    exit
fi

if [[ "$1" -eq "linux32" ]]; then
	make_version
    build_linux32
    exit
fi

if [[ "$1" -eq "linux64" ]]; then
	make_version
    build_linux64
    exit
fi

if [[ "$1" -eq "linuxarm" ]]; then
	make_version
    build_linuxarm
    exit
fi

echo Invalid argument! 
echo Use any of: 
echo ./build.sh 
echo ./build.sh win32 
echo ./build.sh win64 
echo ./build.sh linux32 
echo ./build.sh linux64 
echo ./build.sh linuxarm