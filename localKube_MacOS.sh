#!/bin/bash

curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-darwin-amd64 \
&& sudo install minikube-darwin-amd64 /usr/local/bin/minikube

curl -LO https://storage.googleapis.com/minikube/releases/latest/docker-machine-driver-hyperkit \
&& sudo install -o root -g wheel -m 4755 docker-machine-driver-hyperkit /usr/local/bin/

echo -e "\033[31mplease run minikube start --vm-driver hyperkit\033[0m"

echo -e "\033[31mAfter using the cluster, you can use\033[0m"
echo -e "\033[31mminikube stop\033[0m"
echo -e "\033[31mTo stop the cluster\033[0m"
echo -e "\033[31mminikube delete\033[0m"
echo -e "\033[31mTo delete the cluster\033[0m"
