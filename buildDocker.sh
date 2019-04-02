#!/bin/bash

source config.sh

docker login -u ${userName} -p ${passWd}

# build GAProj request handler
docker build -f Dockerfile.dataservice -t ${gaProjectDockerIamge}:latest .
docker push ${gaProjectDockerIamge}:latest

# build island container
docker build -f Dockerfile.island -t ${islandDockerIamge}:latest .
docker push ${islandDockerIamge}:latest

docker rmi ${gaProjectDockerIamge}:latest
docker rmi ${islandDockerIamge}:latest

docker rmi $(docker images | grep "^<none>" | awk "{print $3}") &> /dev/null