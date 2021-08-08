#!/usr/bin/env bash

set -ex

CLUSTER_NAME=${1:-"default"}

yc managed-kubernetes cluster delete --name k8s-$CLUSTER_NAME
