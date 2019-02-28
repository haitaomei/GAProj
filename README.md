# GAProj

Install CLI Tools (kubetcl, and Helm)
======
1. [kubetcl](https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-with-homebrew-on-macos)
2. [Helm](https://helm.sh/docs/using_helm/#installing-helm)

Install Docker
======
See https://docs.docker.com/install/

Create a Kubernetes clsuter in IBM Kubernetes Service, and install CLI according to the instructions
======
After the cluster online, and setting up the CLI, then type `helm init` to initialise Helm.


Install Project into Kubernetes Cluster
======
Run `install.sh`

Configure Redis Cluster
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

Wait all redis container online, make a note of the name of the first redis container, in this case it's `redis-cluster-788d7c769c-64p45`.

Then type `kubectl exec -it redis-cluster-788d7c769c-64p45 -- redis-cli --cluster create --cluster-replicas 1 $(kubectl get pods -l app=redis-cluster -o jsonpath='{range.items[*]}{.status.podIP}:6379 ')` to set up the redis cluster.

By default, there will be 3 master nodes, and 3 slave nodes. Config `gaproj/values.yaml` to configure the total number of instances.

Exported service url is `redis-cluster.svc.cluster.local:6379`

Delete Project from Kubernetes Cluster
======
Run `cleanup.sh`




Build Docker images (Optional)
======
Copy `config.sh.tpl` to `config.sh`, configure your docker hub's username and password etc.

Move the directory of this project into your $GOPATH/src/

Run `buildDocker.sh`