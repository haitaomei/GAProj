#!/bin/bash

source config.sh

docker login -u ${userName} -p ${passWd}

# docker build -t khitaomei/travis-ci:${TRAVIS_BUILD_NUMBER} .
# docker push khitaomei/travis-ci:${TRAVIS_BUILD_NUMBER}