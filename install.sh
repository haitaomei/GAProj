#!/bin/bash

source config.sh

# generate deployment for islands
i=1
while [ $i -le $numberOfIlands ]
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
    image: {{ .Values.islanddockerimage }}
    env:
    - name: ISLAND_ID
      value: "$i"
EOF

  i=$(( $i + 1 ))
done


ARG=""
ARG="${ARG} --set gaserverdockerimage=${gaProjectDockerIamge}:latest"
ARG="${ARG} --set islanddockerimage=${islandDockerIamge}:latest"
ARG="${ARG} --set gaserverreplica=${gaserverreplica}"
# install the project into kubernetes
helm install --name ga-project ${ARG} gaproject

# delete tmp files
rm -rf gaproject/templates/island*