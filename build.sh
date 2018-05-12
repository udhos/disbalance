#!/bin/bash

go get github.com/gopherjs/gopherjs
go get gopkg.in/yaml.v2

gofmt -s -w ./*/*.go
go tool fix ./*/*.go

go tool vet ./disbalance
go test -v ./disbalance
go install ./disbalance

gopherjs build ./console -o ./run/console/console.js
