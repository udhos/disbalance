#!/bin/bash

gofmt -s -w ./*/*.go
go tool fix ./*/*.go
go tool vet ./disbalance
go test -v ./disbalance
go install ./disbalance

echo openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem
