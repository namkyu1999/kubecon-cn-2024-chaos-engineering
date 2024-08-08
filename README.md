# KubeCon China 2024 Chaos Engineering

## Introduction

## Architecture

### V1

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
### V1

### V2

```shell
## init dapr in k8s
kubectl create namespace v2
dapr init --kubernetes --wait -n v2

## init redis for pubsub
## https://docs.dapr.io/getting-started/tutorials/configure-state-pubsub/#step-1-create-a-redis-store
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install redis bitnami/redis --set image.tag=6.2 --set architecture=standalone -n v2
kubectl apply -f ./v2/deploy/redis-state.yaml

## install rabbitmq for pubsub
helm install rb bitnami/rabbitmq-cluster-operator -n v2
kubectl apply -f ./v2/deploy/rabbitmq.yaml
kubectl apply -f ./v2/deploy/pubsub.yaml

## install mongodb
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install mongo bitnami/mongodb --values ./v2/deploy/mongo-values.yaml -n v2 --version 12.1.31

## install applications
kubectl apply -f ./v2/deploy/delivery.yaml
kubectl apply -f ./v2/deploy/order.yaml
```

# Reference
- [Dapr tutorial - hello-kubernetes](https://github.com/dapr/quickstarts/tree/master/tutorials/hello-kubernetes)