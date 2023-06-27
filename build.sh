#!/bin/bash

GOOS=darwin GOARCH=amd64 go build -o beanjs-v2ray-mac-amd64 .
GOOS=darwin GOARCH=arm64 go build -o beanjs-v2ray-mac-arm64 .

GOOS=linux GOARCH=amd64 go build -o beanjs-v2ray-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -o beanjs-v2ray-linux-arm64 .

GOOS=windows GOARCH=amd64 go build -o beanjs-v2ray-win-amd64.exe .
GOOS=windows GOARCH=arm64 go build -o beanjs-v2ray-win-arm64.exe .