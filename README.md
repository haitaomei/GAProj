# GAProj

Install CLI Tools (kubetcl, and Helm)
------
1. [kubetcl](https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-with-homebrew-on-macos)
2. [Helm](https://helm.sh/docs/using_helm/#installing-helm)

Install Docker
----
See https://docs.docker.com/install/

Create a Kubernetes clsuter in IBM Kubernetes Service, and install CLI according to the instructions
---

Build Docker images
---
Copy `config.sh.tpl` to `config.sh`, configure your docker hub's username and password 

Run `buildDocker.sh`