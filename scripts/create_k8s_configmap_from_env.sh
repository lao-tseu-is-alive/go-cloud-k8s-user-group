#!/bin/bash
# https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/
kubectl create configmap app-config-go-cloud-k8s-user-group --from-env-file=.env --output yaml --dry-run=client