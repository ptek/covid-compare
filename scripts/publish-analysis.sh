#!/bin/bash -l

set -eu
set -o pipefail

s3Key=${S3_KEY}
s3Secret=${S3_SECRET}
file=${PROJECT_ROOT}/data/data-incidences.csv
bucket=670c0cb7634c45da1473a86c9fcafd8c3183e85fcb411849377272c084939

./scripts/s3upload.sh "${s3Key}" "${s3Secret}" ${bucket}@eu-central-1 "${file}" "cases.csv"