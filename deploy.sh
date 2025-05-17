#!/bin/bash

echo "Creating namespace 'movieexample'..."
kubectl create namespace movieexample

echo "Applying ConfigMap and Secret..."
kubectl apply -n movieexample -f config/configmap.yaml
kubectl apply -n movieexample -f config/secret.yaml

echo "Deploying MySQL..."
kubectl apply -n movieexample -f mysql/deployment.yaml
kubectl apply -n movieexample -f mysql/service.yaml

echo "Deploying Kafka..."
kubectl apply -n movieexample -f kafka/deployment.yaml
kubectl apply -n movieexample -f kafka/service.yaml

echo "Deploying Consul..."
kubectl apply -n movieexample -f consul/deployment.yaml
kubectl apply -n movieexample -f consul/service.yaml

echo "Deploying Microservices..."
for svc in rating metadata movie
do
  kubectl apply -n movieexample -f $svc/deployment.yaml
  kubectl apply -n movieexample -f $svc/service.yaml
done

echo "All components deployed in 'movieexample' namespace."
