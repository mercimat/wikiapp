#!/bin/sh
export PATH=$PATH:/usr/local/go/bin
go get -d -v go.mongodb.org/mongo-driver/mongo
go test -v -vet=off . -run TestFrontPageHandler
