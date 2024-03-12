# Building Cloud Native Services with Go + Dapr + Kubernetes
This repository contains examples showing how to build cloud-native services with Go (golang) and Dapr running on Kubernetes.

## What is Dapr?
The Distributed Application Runtime ([Dapr](https://github.com/dapr)) is a Cloud Native Computing Foundation (CNCF) project that exposes a set of tools, APIs, and components for creating distributed applications that can be executed in different environmental contexts including: 

* Self-hosted (local, VMs, or bare metal servers), 
* Kubernetes 
* Or, cloud provider managed

The examples in this repository assume they are running in a Kubernetes cluster (local or otherwise).

## Run examples
The examples in this repository are designed to showcase the diverse components and building blocks of Dapr. To run the examples you will need a local Kubernetes cluster running KinD or Minikube. Find the complete list of requirements below.

### Requirements
Local requirements to run examples:

* Dapr CLI
* Docker or similar tool 
* Kubernetes cluster (on KinD or Minikube)
* Dapr control plane on Kubernetes
* Latest Go version 
* Ko container image build tool

### Install the Dapr CLI
One of the first thing to do is to install the Docker CLI.

* Use [these instructions](https://docs.dapr.io/getting-started/install-dapr-cli/) to install the Dapr CLI

### Create a local clsuter

You can  use the KinD tool (minikuber or similar tools) to launch a local cluster. The examples in this repository assume a Kubernetres cluster running the Dapr runtime.

* Follow instructions to [create a Kind cluster](https://kind.sigs.k8s.io/docs/user/quick-start/)
* Or, follow instructions to [create a Minikube cluster](https://minikube.sigs.k8s.io/docs/start/)

### Deploy Dapr control plane
Once you have a Kubernetes cluster, the next step is to deploy the Dapr control plane on the cluster.

The easiest way to do this is to use the Dapr CLI tool to deploy the Dapr Kubernetes components to the cluster.

```
dapr init --kubernetes
```

The previous command will deploy Dapr controllers and other control plane components to run Dapr services on Kubernetes. For additional detail, see [Deploy Dapr on Kubernetes](https://docs.dapr.io/operations/hosting/kubernetes/kubernetes-deploy/).

### Install Go
You will need the Go build tools to compile the example code in this repository.

* Find and [install the latest Go tools](https://go.dev/doc/install)

### Install Ko
Ko is a container image builder for Go applications. 

* [Install ko](https://ko.build/install/) on your machine

### Examples 

* 01-service-invoke - Shows how to setup a simple service invocation