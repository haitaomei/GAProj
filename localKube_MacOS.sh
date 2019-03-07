# !/bin/bash

curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-darwin-amd64 \
&& sudo install minikube-darwin-amd64 /usr/local/bin/minikube

curl -LO https://storage.googleapis.com/minikube/releases/latest/docker-machine-driver-hyperkit \
&& sudo install -o root -g wheel -m 4755 docker-machine-driver-hyperkit /usr/local/bin/

echo -e "\033[31mplease run\033[0m"
echo "minikube start --vm-driver hyperkit"

echo -e "\033[31mAfter minikube is running, using the following commands to enable ingree\033[0m"
echo "minikube addons enable ingress"

echo -e "\033[31mIf you want to enable HPA in minikube, using the following command\033[0m"
echo "minikube addons enable heapster"
echo "minikube addons enable metrics-server"


echo -e "\033[31mAfter using the cluster, you can use\033[0m"
echo "mminikube stop"
echo -e "\033[31mTo stop the cluster\033[0m"
echo "minikube delete"
echo -e "\033[31mTo delete the cluster\033[0m"

