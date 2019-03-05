#!/bin/bash

source config.sh

# generate deployment for islands
for i in {1..2}
do
    cat <<EOF >>gaproject/templates/island$i.yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: ga-island-$i
spec:
  containers:
  - name: ga-island-$i
    image: redis
    env:
    - name: ISLAND_ID
      value: "$i"
EOF
done

# install the project into kubernetes
helm install --name ga-project gaproject

# delete tmp files
rm -rf gaproject/templates/island*