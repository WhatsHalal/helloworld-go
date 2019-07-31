#!/bin/bash

VERSION=$(git rev-parse HEAD)
mkdir -p bin
go build -o bin/$VERSION/helloworld .
tar zcf bin/helloworld_$VERSION.tar.gz -C bin/$VERSION .

