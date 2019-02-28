#!/bin/bash

source config.sh

docker login -u ${userName} -p ${passWd}

# build GAProj request handler
rm -rf GAProj
GOOS=linux go build
docker build -f Dockerfile.gaproj -t ${gaProjectDockerIamge}:latest .
docker push ${gaProjectDockerIamge}:latest

rm -rf GAProj