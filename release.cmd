@echo off
echo Building binary ...

go build -ldflags -H=windowsgui

echo Copying to Utils ...
copy ssh-chooser.exe c:\utils\