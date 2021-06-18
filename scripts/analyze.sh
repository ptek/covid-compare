#!/bin/bash

set -eu
set -o pipefail

docker exec -i spark-master mkdir -p /app/
docker cp ../target/scala-2.11/covid-compare_2.11-0.1.jar spark-master:/app/
docker exec -i spark-master /spark/bin/spark-submit --class CovidDataComparison /app/covid-compare_2.11-0.1.jar
