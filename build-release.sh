#!/bin/bash

# BUILD

# Get Go version
GO_VERSION=$(go version | awk '{print $3}')

# Get the build date
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

go build -o discord-bot-boilerplate -ldflags "-s -X github.com/keshon/discord-bot-boilerplate/internal/version.BuildDate=$BUILD_DATE -X github.com/keshon/discord-bot-boilerplate/internal/version.GoVersion=$GO_VERSION" cmd/main.go

upx discord-bot-boilerplate