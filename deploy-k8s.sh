#!/usr/bin/env bash

set -xe

CLUSTER_NAME=${1:-"default"}
CONFIG_PATH=$2

NODES=`cat $CONFIG_PATH | jq -r .nodes`
CPU=`cat $CONFIG_PATH | jq -r .cpu`
RAM=`cat $CONFIG_PATH | jq -r .ram`

FOLDER_ID=$(yc config get folder-id)
RES_SA_ID=$(yc iam service-account get --name k8s-res-sa-${FOLDER_ID} --format json | jq .id -r)

yc managed-kubernetes cluster create \
--name k8s-$CLUSTER_NAME \
--network-name default \
--zone ru-central1-a \
 --subnet-name default-ru-central1-a \
--public-ip \
--service-account-id $RES_SA_ID \
--node-service-account-id $RES_SA_ID

yc managed-kubernetes node-group create \
--name k8s-$CLUSTER_NAME-node-group \
--cluster-name k8s-$CLUSTER_NAME \
--platform-id standard-v2 \
--public-ip \
--cores $CPU \
--memory $RAM \
--core-fraction 50 \
--disk-type network-ssd \
--fixed-size $NODES \
--location subnet-name=default-ru-central1-a,zone=ru-central1-a

yc managed-kubernetes cluster get-credentials --external --name k8s-$CLUSTER_NAME --force

CLUSTER_IP=`kubectl get nodes -o json | jq -r '.items[0].status.addresses[] | select(.type=="ExternalIP") | .address'`
echo $CLUSTER_IP > cluster_ip.txt
