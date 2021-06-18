#!/bin/bash -l

set -eux
set -o pipefail

aws s3 cp output/cases.csv s3://670c0cb7634c45da1473a86c9fcafd8c3183e85fcb411849377272c084939/cases.csv --acl public-read