#!/bin/bash

set -eu
set -o pipefail

get-de() {
  result="${PROJECT_ROOT}/data/data-de.csv"
  mkdir -p "$(dirname $result)"
  wget -nv -O "$result" "https://arcgis.com/sharing/rest/content/items/f10774f1c63e40168479a1feb6c7ca74/data"
}

get-pl() {
  mkdir -p "${PROJECT_ROOT}/data"
  result="${PROJECT_ROOT}/data/data-pl.csv"
  archive="${PROJECT_ROOT}/data/data-pl.zip"
  archive_extract=$(mktemp -d)

  wget -nv -O "$archive" "https://arcgis.com/sharing/rest/content/items/e16df1fa98c2452783ec10b0aea4b341/data"
  unzip -q -o "$archive" -d "$archive_extract"

  find "$archive_extract" -name "*_rap_rcb_pow_eksport.csv" -print0 | xargs -n1 -0 -I{} enconv -L pl -x UTF-8 {}
  find "$archive_extract" -name "*_rap_rcb_pow_eksport.csv" -print0 | xargs -n1 -0 -I{} bash -c "dt=\$(basename {} | head -c 8); sed -i -e \"s/^/\$dt;/\" {}"

  echo -n > $result
  find "$archive_extract" -name "*_rap_rcb_pow_eksport.csv" -print0 | xargs -n1 -0 -I{} bash -c "tail +2 {} >>$result"

  rm -rf "$archive_extract"
}

main() {
  get-de
  get-pl
}
main