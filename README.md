# GAProj

1. Install CLI Tools (kubetcl, and Helm)
======
* [kubetcl](https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-with-homebrew-on-macos)
* [Helm](https://helm.sh/docs/using_helm/#installing-helm)


2. Install Docker
======
See https://docs.docker.com/install/

Create a Kubernetes clsuter 
======
Create in IBM Kubernetes Service, and install CLI according to the instructions.

Or if you want to run locally on MacOS, please install and start minikube using `localKube_MacOS.sh`

**After** the cluster online, and setting up the CLI, then type `helm init` to initialise Helm.

3. Config
======
Copy `config.sh.tpl` to `config.sh`, configure your docker hub's username and password etc.

4. (Optional) Build Docker images 
======
Run `buildDocker.sh`


5. Install Project into Kubernetes Cluster
======
Run `install.sh`

For now we use a single redis instance. If we set up Redis cluster on Kubernetes without persistent storage, it makes no sense.

<!-- Configure Redis Cluster
----
Then type `kubectl get pod`, you will see:

        NAME                             READY     STATUS              RESTARTS   AGE
        redis-cluster-788d7c769c-64p45   0/1       Running             0          26s
        redis-cluster-788d7c769c-bp7ql   0/1       ContainerCreating   0          26s
        redis-cluster-788d7c769c-d7hqn   0/1       Running             0          26s
        redis-cluster-788d7c769c-jr8rm   0/1       ContainerCreating   0          26s
        redis-cluster-788d7c769c-rgzsp   0/1       Running             0          26s
        redis-cluster-788d7c769c-x6fq6   0/1       ContainerCreating   0          26s
        ......

Wait all redis container online, then type 

`kubectl exec -it $(kubectl get pods | grep -m1 redis | awk ' { print $1 } ') -- redis-cli --cluster create --cluster-replicas 1 $(kubectl get pods -l app=redis-cluster -o jsonpath='{range.items[*]}{.status.podIP}:6379 ')` 

to set up the redis cluster.

By default, there will be 3 master nodes, and 3 slave nodes.

Exported service url is `redis-cluster.default.svc.cluster.local:6379` -->

Delete Project from Kubernetes Cluster
======
Run `cleanup.sh`
