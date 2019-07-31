#!/bin/bash

VERSION=$(git rev-parse HEAD)
aws s3 cp bin/helloworld_$VERSION.tar.gz s3://circleci-demo.whatshalal.com/builds/helloworld_$VERSION.tar.gz

