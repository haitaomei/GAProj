#!/bin/bash

source config.sh

docker login -u ${userName} -p ${passWd}

rm -rf GAProj
GOOS=linux go build

docker build -t ${gaProjectDockerIamge}:latest .
docker push ${gaProjectDockerIamge}:latest

rm -rf GAProj