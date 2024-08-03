# KubeCon China 2024 Chaos Engineering

## Introduction
## Setup

```shell
## init dapr in k8s
dapr init --kubernetes --wait -n dapr

## init redis for pubsub
## https://docs.dapr.io/getting-started/tutorials/configure-state-pubsub/#step-1-create-a-redis-store
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install redis bitnami/redis --set image.tag=6.2 --set architecture=standalone -n dapr
kubectl apply -f ./v2/deploy/redis-state.yaml

## install rabbitmq for pubsub
helm install rb bitnami/rabbitmq-cluster-operator -n dapr
kubectl apply -f ./v2/deploy/rabbitmq.yaml
kubectl apply -f ./v2/deploy/pubsub.yaml

## install mongodb
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install mongo bitnami/mongodb --values ./v2/deploy/mongo-values.yaml -n dapr --version 12.1.31

## install applications
kubectl apply -f ./v2/deploy/delivery.yaml
kubectl apply -f ./v2/deploy/order.yaml
```

# Reference
- [Dapr tutorial - hello-kubernetes](https://github.com/dapr/quickstarts/tree/master/tutorials/hello-kubernetes)