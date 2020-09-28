#!/bin/sh
export PATH=$PATH:/usr/local/go/bin
export GOPATH="$WORKSPACE"

mkdir -p $GOPATH/src/github.com/mercimat
ln -s $WORKSPACE $GOPATH/src/github.com/mercimat/wikiapp

go get -d -v go.mongodb.org/mongo-driver/mongo

cd web/
go test -v -vet=off .
