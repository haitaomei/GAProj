#!/bin/bash

source config.sh

docker login -u ${userName} -p ${passWd}

# build GAProj request handler
rm -rf GAProj
GOOS=linux go build
docker build -f Dockerfile.gaproj -t ${gaProjectDockerIamge}:latest .
docker push ${gaProjectDockerIamge}:latest

# build island container
docker build -f Dockerfile.island -t ${islandDockerIamge}:latest .
docker push ${islandDockerIamge}:latest

docker rmi ${gaProjectDockerIamge}:latest
docker rmi ${islandDockerIamge}:latest

rm -rf GAProj