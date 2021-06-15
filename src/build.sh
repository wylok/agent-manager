#!/usr/bin/env bash
go build -ldflags "-w -s" main.go
mv -f main manager
upx --brute manager