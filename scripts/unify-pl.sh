#!/bin/bash

set -eu
set -o pipefail

realpath_osx() {
    path=$(eval echo "$1")
    folder=$(dirname "$path")
    res=$(cd "$folder"; pwd)/$(basename "$path")
    echo "$res"
}

SCRIPT=$(realpath_osx "$0")
SCRIPTPATH=$(dirname "$SCRIPT")

pushd "$SCRIPTPATH"
pushd raw-data/data-pl/

find . -name "*_rap_rcb_pow_eksport.csv" -print0 | tail -n1 | xargs head -n1 >./pl.csv
# find . -name "*_rap_rcb_pow_eksport.csv" -print0 | xargs -L1 -0 -I{} tail +2 {} | iconv -f WINDOWS-1250 --unicode-subst=. --byte-subst=. --widechar-subst=. -t UTF-8 >>./pl.csv
find . -name "*_rap_rcb_pow_eksport.csv" -print0 | xargs -L1 -0 -I{} bash -c "tail +2 {} | enconv -L pl -x UTF-8 >>./pl.csv"

popd
popd