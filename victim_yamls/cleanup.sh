#!/bin/bash

BASE_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

yamls=($(ls $BASE_DIR/*.yaml))

for victim_yaml in "${yamls[@]}"; do
  kubectl delete -f $victim_yaml || true
done
