# KubeCon China 2024 Chaos Engineering

## Introduction
- Session Info: [What if Your System Experiences an Outage? Let's Build a Resilient Systems with Chaos Engineering](https://sched.co/1eYaZ)
- Recording: TBD
- Slide: TBD

## Architecture

### V1
![v1](./img/v1.png)
### V2
![v2](./img/v2.png)

### Chaos Engineering Info

![duration](img/observations.png)

## Setup

### Pre-requisite

```shell
## start minikube
minikube start --memory 8192 --cpus 4 

## init dapr in k8s
dapr init --kubernetes --wait -n default

## install rabbitmq operator for pubsub
helm install rb bitnami/rabbitmq-cluster-operator -n default

## install prometheus using helm
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install dapr-prom prometheus-community/prometheus --values ./common/prometheus-values.yaml -n default

## install Grafana
## username: admin, password: admin
helm repo add grafana https://grafana.github.io/helm-charts
kubectl --namespace default create secret generic grafana-password \
   --from-literal=admin-user=admin --from-literal=admin-password=admin
helm install grafana grafana/grafana --namespace default \
  --set admin.existingSecret=grafana-password

```

### LitmusChaos

```shell
## install mongo
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install my-release bitnami/mongodb --values ./litmus/mongodb-values.yaml -n litmus --create-namespace --version 12.1.31

## install litmus v3.8.0
kubectl apply -f https://raw.githubusercontent.com/litmuschaos/litmus/master/mkdocs/docs/3.8.0/litmus-cluster-scope-3.8.0.yaml

kubectl create secret generic k6-script-v1 \
    --from-file=./litmus/script-v1.js -n litmus
    
kubectl create secret generic k6-script-v2 \
    --from-file=./litmus/script-v2.js -n litmus
```

### V1
```shell
kubectl create namespace v1
kubectl apply -f ./v1/deploy/rabbitmq.yaml
kubectl apply -f ./v1/deploy/pubsub.yaml
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

## install rabbitmq for pubsub
kubectl apply -f ./v2/deploy/rabbitmq.yaml
kubectl apply -f ./v2/deploy/pubsub.yaml

## install applications
kubectl apply -f ./v2/deploy/delivery.yaml
kubectl apply -f ./v2/deploy/order.yaml
```

# Reference
- [Dapr - hello-kubernetes](https://github.com/dapr/quickstarts/tree/master/tutorials/hello-kubernetes)
- [Dapr - how to use outbox](https://docs.dapr.io/developing-applications/building-blocks/state-management/howto-outbox/)
- [Dapr - rabbitmq](https://docs.dapr.io/reference/components-reference/supported-pubsub/setup-rabbitmq/)
- [Dapr - Prometheus Setup](https://docs.dapr.io/operations/observability/metrics/prometheus/)
- [RabbitMQ's persistence](https://www.rabbitmq.com/kubernetes/operator/using-operator#persistence)
- [RabbitMQ HA mode](https://www.infracloud.io/blogs/setup-rabbitmq-ha-mode-kubernetes-operator/)
- [MongoDB cli auth](https://medium.com/@yasiru.13/mongodb-setting-up-an-admin-and-login-as-admin-856ea6856faf)
- [RabbitMQ Prometheus Setup](https://www.rabbitmq.com/kubernetes/operator/operator-monitoring)
- [RabbitMQ Dashboard](https://grafana.com/grafana/dashboards/10991-rabbitmq-overview/)