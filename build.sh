#!/bin/sh

export GOPROXY=https://goproxy.cn,direct
export env GOOS=linux
export GOARCH=arm
export CGO_ENABLED=1
export GOARM=7
export CC=arm-linux-gnueabihf-musl-gcc

# libffi.so.8
export LIBFFI_PATH=/home/ghost/software/libffi-3.4.2/_install
# libffi.so.7
# export LIBFFI_PATH=/home/ghost/software/libffi-3.3/_install

export CGO_CFLAGS="-I$LIBFFI_PATH/include"
export CGO_LDFLAGS="-L$LIBFFI_PATH/lib"

rm -rf dist

go build -p 12 -ldflags "-s -w" -o dist/usg