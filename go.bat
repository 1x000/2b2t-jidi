@echo off

set GOOS=linux
set GOARCH=arm64
set CGO_ENABLED=0

go build -o arm64linux

if %errorlevel% equ 0 (
    echo Compilation successful: arm64linux
) else (
    echo Compilation failed for %GOOS%/%GOARCH%
    exit /b %errorlevel%
)
