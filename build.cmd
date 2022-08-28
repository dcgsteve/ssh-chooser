@echo off
echo Building binary ...

go-winres make
env GO111MODULE=on go build -ldflags "-H=windowsgui"
