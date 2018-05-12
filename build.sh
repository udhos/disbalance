#!/bin/bash

go get github.com/gopherjs/gopherjs
go get gopkg.in/yaml.v2

gofmt -s -w ./*/*.go
go tool fix ./*/*.go

go tool vet ./disbalance
go test -v ./disbalance
go install ./disbalance

go build -o ./run/disbalance ./disbalance

gopherjs build -o ./run/console/console.js ./console
