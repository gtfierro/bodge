#!/bin/bash

set -e

if (( "$#" >= 1 )); then
    export BW2_DEFAULT_ENTITY=$1
    echo "Using $BW2_DEFAULT_ENTITY"
else if [[ -v BW2_DEFAULT_ENTITY ]]; then
    echo "Using $BW2_DEFAULT_ENTITY"
else
    echo 'Please either supply an entity file or set $BW2_DEFAULT_ENTITY'
    exit 1
fi
fi

go generate -tags ragent
go build -tags ragent
