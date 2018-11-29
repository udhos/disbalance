#!/bin/bash

go get github.com/gopherjs/gopherjs

#go get gopkg.in/yaml.v2
#go get honnef.co/go/js/dom

gofmt -s -w ./*/*.go
go tool fix ./*/*.go
go tool vet ./console
go tool vet ./disbalance

hash golint && golint disbalance rule console

go test ./disbalance
go install ./disbalance

go build -o ./run/disbalance ./disbalance

gopherjs build -o ./run/console/console.js ./console
