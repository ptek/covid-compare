#!/bin/bash

set -ex

if [ -z "$SKIP_DOWNLOAD" ]
then
	./get-data.sh
fi
./upload-data.sh
./analyze.sh 2> /dev/null
./download-analysis.sh
./publish-analysis.sh
