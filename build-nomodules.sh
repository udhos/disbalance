#!/bin/bash

go get github.com/gopherjs/gopherjs

go get gopkg.in/yaml.v2
go get honnef.co/go/js/dom

gofmt -s -w ./*/*.go
go tool fix ./*/*.go
go vet ./console
go vet ./disbalance

hash golint 2>/dev/null && golint disbalance rule console

go test ./disbalance
go install ./disbalance

go build -o ./run/disbalance ./disbalance

gopherjs build -o ./run/console/console.js ./console
