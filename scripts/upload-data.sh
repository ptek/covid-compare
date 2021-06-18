#!/bin/bash

set -eu
set -o pipefail

docker cp raw-data namenode:/tmp/
docker exec -i namenode hdfs dfs -rm -f -r /user/root/input/raw-data
docker exec -i namenode hdfs dfs -copyFromLocal /tmp/raw-data/ /user/root/input/