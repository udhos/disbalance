#!/bin/bash

export CGO_ENABLED=0
export GO111MODULE=on

gofmt -s -w ./*/*.go
go tool fix ./*/*.go
go vet ./console
go vet ./disbalance

hash golint 2>/dev/null && golint disbalance rule console

go test ./disbalance
go install ./disbalance

go build -o ./run/disbalance ./disbalance

gopherjs build -o ./run/console/console.js ./console
