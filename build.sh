#!/bin/bash

go get github.com/gopherjs/gopherjs

gofmt -s -w ./*/*.go
go tool fix ./*/*.go

go tool vet ./disbalance
go test -v ./disbalance
go install ./disbalance

gopherjs build ./disbalance -o ./run/console/console.js
