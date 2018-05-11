#!/bin/bash

gofmt -s -w ./*/*.go
go tool fix ./*/*.go
go tool vet ./disbalance
go test -v ./disbalance
go install ./disbalance

