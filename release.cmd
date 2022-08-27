@echo off
echo Building binary ...

go build -ldflags -H=windowsgui
