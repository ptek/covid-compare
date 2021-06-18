#!/bin/bash

set -eu
set -o pipefail

docker exec -i namenode rm -rf /tmp/cases.csv
docker exec -i namenode hadoop fs -getmerge /user/root/output/cases.csv /tmp/cases.csv
docker cp namenode:/tmp/cases.csv output/cases.csv