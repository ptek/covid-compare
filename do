#!/bin/bash

set -eu
set -o pipefail

PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
export PROJECT_ROOT

build() {
    rm -rf dist
    mkdir -p "${PROJECT_ROOT}/dist/bin"
    
    pushd ./src
        go mod tidy
        env GOOS=linux go build -o "${PROJECT_ROOT}/dist/bin/covid-compare_Linux"
        go build -o "${PROJECT_ROOT}/dist/bin/covid-compare_Darwin"
    popd

}

container_name="ghcr.io/ptek/covid-compare"
build_container() {
    podman build -t ${container_name} -f "${PROJECT_ROOT}/docker/Dockerfile" .
}


push_container() {
    podman push ${container_name}
}

run_dev() {
    pushd src
        go run main.go
    pops
}

run() {
    if [ "${SKIP_DOWNLOAD:="no"}" == "no" ]
    then
        ./scripts/get-data.sh
    fi
    ./dist/bin/covid-compare_$(uname)
    ./scripts/publish-analysis.sh
}

usage() {
    echo "./do <command>"
    echo ""
    echo "command can be one of:"
    echo "  build:   build the binaries for the local system and linux x86_64"
    echo "           and put them into the dist folder"
    echo ""
    echo "  build_container: build an OCI container with the scripts needed"
    echo "           to run all parts of the analysis"
    echo ""
    echo "  run_dev: compile and run the latest version of the code locally"
    echo ""
    echo "  run:     run the scripts to fetch the data, run analysis, and publish"
    echo "           the results"
    echo ""
}

main() {
  if [ -z "$@" ]
  then
    usage
  else
    for arg in "$@"
    do
        case "$arg" in
            "build" )
                build;;
            "build_container" )
                build_container;;
            "push_container" )
                push_container;;
            "run_dev" )
                run_dev;;
            "run" )
                run;;
            * )
                usage;;
        esac
    done
  fi
}

main $@