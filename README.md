# KubeCon China 2024 Chaos Engineering

## Introduction

## Architecture

### V1
![v1](./img/v1.png)
### V2
![v2](./img/v2.png)

### Chaos Engineering Info

![duration](img/observations.png)

## Setup

### LitmusChaos

```shell
## install mongo
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install my-release bitnami/mongodb --values ./litmus/mongodb-values.yaml -n litmus --create-namespace --version 12.1.31

## install litmus v3.8.0
kubectl apply -f https://raw.githubusercontent.com/litmuschaos/litmus/master/mkdocs/docs/3.8.0/litmus-cluster-scope-3.8.0.yaml

kubectl create secret generic k6-script \
    --from-file=./litmus/script.js -n litmus
```

### Common
```shell
kubectl create namespace common
dapr init --kubernetes --wait -n common

## install mongodb
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install mongo bitnami/mongodb --values ./common/mongo-values.yaml -n common --version 12.1.31

## install rabbitmq for pubsub
helm install rb bitnami/rabbitmq-cluster-operator -n common
kubectl apply -f ./common/rabbitmq.yaml
kubectl apply -f ./v2/deploy/pubsub.yaml
```

### V1
```shell
kubectl create namespace v1
kubectl apply -f ./v1/deploy/delivery.yaml
kubectl apply -f ./v1/deploy/order.yaml
```

### V2

```shell
## init dapr in k8s
kubectl create namespace v2

## init redis for pubsub
## https://docs.dapr.io/getting-started/tutorials/configure-state-pubsub/#step-1-create-a-redis-store
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install redis bitnami/redis --set image.tag=6.2 --set architecture=standalone -n v2
kubectl apply -f ./v2/deploy/redis-state.yaml

## install applications
kubectl apply -f ./v2/deploy/delivery.yaml
kubectl apply -f ./v2/deploy/order.yaml
```

# Reference
- [Dapr tutorial - hello-kubernetes](https://github.com/dapr/quickstarts/tree/master/tutorials/hello-kubernetes)