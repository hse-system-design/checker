#!/usr/bin/env bash

set -xe

TANK_NAME=${1:-default}

yc compute instances delete --name yandex-tank-$TANK_NAME
