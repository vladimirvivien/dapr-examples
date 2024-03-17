# Building Cloud Native Services with Go + Dapr + Kubernetes
This repository contains examples showing how to build cloud-native services with Go (golang) and Dapr running on Kubernetes.

## What is Dapr?
The Distributed Application Runtime ([Dapr](https://github.com/dapr)) is a Cloud Native Computing Foundation (CNCF) project that exposes a set of tools, APIs, and components for creating distributed applications that can be executed in different environmental contexts including: 

* Self-hosted (local, VMs, or bare metal servers), 
* Kubernetes 
* Or, cloud provider managed

The examples in this repository assume they are running in a Kubernetes cluster (local or otherwise).

## Environment setup
The Dapr examples use several tools to work properly on your local machine. See the requirements below and ensure your machine is setup with all the necessary tools prior to trying the examples.

### Requirements
Local requirements to run examples:

* Dapr CLI
* Docker or similar tool 
* Kubernetes cluster (on KinD or Minikube)
* Dapr control plane on Kubernetes
* Latest Go version 
* Ko container image build tool
* Redis

### Install the Dapr CLI
One of the first thing to do is to install the Docker CLI.

* Use [these instructions](https://docs.dapr.io/getting-started/install-dapr-cli/) to install the Dapr CLI

### Install Docker
Your environment will need Docker or similar tool to host your OCI-compliant image and as a container runtime.

* See instructions on [Install Docker Engine](https://docs.docker.com/engine/install/)

### Install Kind Kubernetes cluster management tool

You can  use the KinD tool (minikuber or similar tools) to launch a local cluster. The examples in this repository assume a Kubernetres cluster running the Dapr runtime.

* Follow instructions to [create a Kind cluster](https://kind.sigs.k8s.io/docs/user/quick-start/)

For instance, given the following file `config/kind-cluster.yaml`:

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  kubeadmConfigPatches:
  - |
    kind: InitConfiguration
    nodeRegistration:
      kubeletExtraArgs:
        node-labels: "ingress-ready=true"    
  extraPortMappings:
  - containerPort: 80
    hostPort: 8081
    protocol: TCP
  - containerPort: 443
    hostPort: 8443
    protocol: TCP
- role: worker
- role: worker
```

Then, the following command creates a cluster name `dapr-cluster`:

```
kind create cluster --config ./configs/kind-cluster.yaml --name dapr-cluster
```

### Deploy Dapr control plane
Once you have a Kubernetes cluster, the next step is to deploy the Dapr control plane on the cluster.

The easiest way to do this is to use the Dapr CLI tool to deploy the Dapr Kubernetes components to the cluster.

```
dapr init --kubernetes
```

The previous command will deploy Dapr controllers and other control plane components to run Dapr services on Kubernetes.
Next, you can use the `dapr` command to verify the deployment of the Dapr control plane components on the cluster:

```
dapr status -k

NAME                   NAMESPACE    HEALTHY  STATUS   REPLICAS  VERSION  AGE  CREATED   
dapr-sidecar-injector  dapr-system  True     Running  1         1.13.0   31m  2024-03-16
dapr-sentry            dapr-system  True     Running  1         1.13.0   31m  2024-03-16
dapr-placement-server  dapr-system  True     Running  1         1.13.0   31m  2024-03-16
dapr-operator          dapr-system  True     Running  1         1.13.0   31m  2024-03-16
dapr-dashboard         dapr-system  True     Running  1         0.14.0   31m  2024-03-16
```

For additional detail, see [Deploy Dapr on Kubernetes](https://docs.dapr.io/operations/hosting/kubernetes/kubernetes-deploy/).

### Install Go
You will need the Go build tools to compile the example code in this repository.

* Find and [install the latest Go tools](https://go.dev/doc/install)

### Install Ko
Ko is a container image builder for Go applications. 

* [Install ko](https://ko.build/install/) on your machine

### Install Helm
You will need heml to install certain components used in the examples. 

* Install [helm](https://helm.sh/docs/intro/install/) locally.

### Install Redis
The examples use Redis for both Dapr storage management component and publish/subscribe component. Note that Dapr makes it easy to easily swap out Redis for your preferred data platform as the backing for these components.  See [Dapr Components](https://docs.dapr.io/operations/components/) for detail.

Use helm to install Redis locally.

```
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install redis bitnami/redis
```

## Building the source and loading Docker images
To simplify the local setup, the exampels rely on the `ko` command-line tool to build and generate `Docker` images from the Go source code.

### Compile with ko
You can use `Ko` to compile the source code and automatically push the built image into your local Docker repository. For instance, the following builds code found in package `frontendsvc`:


```
# cd <source code directory>
ko build --local -B --platform=linux/arm64 frontendsvc/front.go
```
The previous step will build and publish an images to the local repository. Notice the optional `--platform` argument to specify the platform for the image to build. 

```
docker images

REPOSITORY             TAG              IMAGE ID       CREATED         SIZE
ko.local/frontendsvc   latest           6094bfc88ad3   2 days ago      16.9MB
```

### Load image into Kind
This step loads the built Docker image into Kind's internal image repository:

```
kind load docker-image ko.local/frontendsvc:latest --name dapr-cluster
```

As an optional step, you can verify that the image is loaded into Kind's internal repository as follows (assuming dapr-cluster as a Kind cluster name):

```
docker exec -it dapr-cluster-control-plane crictl images

IMAGE                    TAG                  IMAGE ID            SIZE
...
ko.local/frontendsvc    latest               1a4f9a427a625       17.4MB
....
```

## Running examples
The examples in this repository are designed to showcase the diverse components and building blocks of Dapr. Each example directory includes a `manifest` directory that contains the YAML configuration for Kubernetes components and services.

### Deploying example service
One of the first steps to run the service is to deploy it unto the Kubernetes cluster.

```
# cd <into example directory>

kubectl apply -f ./manifest
```

### Verify service deployments 
You can use the `kubectl` command to verify the deployment of the different components included in the examples.

First, ensure the Dapr components are deployed properly in the cluster:

```
kubectl get components

NAME           AGE
orders-store   72m
```

Ensure services and deployments are available on the cluster:

```
kubectl get deployments -l app=frontendsvc -o wide

NAME          READY   UP-TO-DATE   AVAILABLE   AGE   CONTAINERS    IMAGES                        SELECTOR
frontendsvc   1/1     1            1           75m   frontendsvc   ko.local/frontendsvc:latest   app=frontendsvc
```

## Troubleshooting
Deploying Dapr-backed services can be complex and may not work the first time.  To save you time, I have gathered some of the troubleshooting steps that can help you identify issues with deploying Dapr services on Kubernetes.

### Code update not reflected
Sometimes your code changes may not be reflected in the cluster. There can be many things that causes this issue:

* Ensure code compiles and is published to the proper container image repository (public or remote)
* Use a version tag (instead of just latest) in the manifest YAML
* If using Kind, load image into the local Kind cluster, then verify correct version is loaded in kind
* If all else fails, delete images from local repository, recompile, and republish

### Dapr sidecar not injected
During initial setup of the many tools needed, it's possible that the Dapr components are installed after you've deployed your causing the Dapr sidecar not to be injected properly.

After your deployment is completed, ensure the Dapr sidecar is being injected into the pod for your pod:

```
kubectl logs -l app=dapr-sidecar-injector -n dapr-system

...
time="2024-03-17T11:58:07.70451326Z" level=info msg="Sidecar injector succeeded injection for app 'frontendsvc'" instance=dapr-sidecar-injector-cb9768b5d-pkl57 scope=dapr.injector.service type=log ver=1.13.0
```

Alternatively, you can describe the pod for your application and make sure the `daprd` container is being injected as a sidecar:

```
kubectl describe pods -l app=frontendsvc
```

The previous command should describe information for pod `frontendsvc` including the injected `daprd` container:

```
Name:             frontendsvc-7c6bb8bf87-znq2d
Namespace:        default
Labels:           app=frontendsvc
                  dapr.io/app-id=frontendsvc
                  dapr.io/metrics-enabled=true
                  dapr.io/sidecar-injected=true
Annotations:      dapr.io/app-id: frontendsvc
                  dapr.io/enabled: true
Status:           Running
Containers:
  frontendsvc:
    Container ID:   containerd://968315781ff6c8c50e336fec2f2ada8413b4820345e08d09f6d7a44e11789080
    Image:          ko.local/frontendsvc:latest
...
  daprd:
    Container ID:  containerd://5eb85acfc539f4c2208318f56c53a57951dae4968426b3a174e64438cf5bb5a1
    Image:         ghcr.io/dapr/daprd:1.13.0
    Args:
      /daprd
...
```

## Examples

* 01-simple-service - A simple HTTP service that saves data into a Dapr-managed data store