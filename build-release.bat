@echo off

rem
rem BUILD
rem

rem Get Go version
for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i

rem Get the build date
for /f "tokens=*" %%a in ('powershell -command "Get-Date -UFormat '%%Y-%%m-%%dT%%H:%%M:%%SZ'"') do set BUILD_DATE=%%a

go build -o discord-bot-boilerplate.exe -ldflags "-s -X github.com/keshon/discord-bot-boilerplate/internal/version.BuildDate=%BUILD_DATE% -X github.com/keshon/discord-bot-boilerplate/internal/version.GoVersion=%GO_VERSION%" cmd\main.go

upx discord-bot-boilerplate.exe