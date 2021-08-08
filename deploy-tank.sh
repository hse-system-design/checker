#!/usr/bin/env bash

set -xe

TANK_NAME=${1:-default}
KEY_PATH=$2

yc compute instance create \
--name yandex-tank-$TANK_NAME \
--network-interface subnet-name=default-ru-central1-a,nat-ip-version=ipv4 \
--zone ru-central1-a \
--ssh-key $KEY_PATH \
--create-boot-disk image-folder-id=standard-images,image-family=yandextank

TANK_IP=`yc compute instance list --format json | jq -r ".[] | select(.name==\"yandex-tank-$TANK_NAME\") | .network_interfaces[0].primary_v4_address.one_to_one_nat.address"`
echo $TANK_IP > tank_ip.txt
