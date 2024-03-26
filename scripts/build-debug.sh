#!/bin/bash

# BUILD

# Get Go version
GO_VERSION=$(go version | awk '{print $3}')

# Get the build date
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Navigate to the root of the project from the scripts folder
cd ..

# Build command
go build -o dbt -ldflags "-X github.com/keshon/discord-bot-template/internal/version.BuildDate=$BUILD_DATE -X github.com/keshon/discord-bot-template/internal/version.GoVersion=$GO_VERSION" cmd/dbt/dbt.go

# Return to the scripts folder after execution
cd scripts