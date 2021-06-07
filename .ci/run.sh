#!/bin/sh

set -e

go test -p=1 -count=1 -v -covermode=count -coverprofile=profile.cov $(go list ./... | grep -v '/mocks/')
go tool cover -func=profile.cov
