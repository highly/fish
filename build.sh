#!/bin/bash

if [[ "$1" == "osx" ]]; then
    os='Darwin'
elif [[ "$1" == "linux" ]]; then
    os='Linux'
fi
if [[ "$os" == "" ]]; then
    os=`uname -s`
fi
if [[ "$os" == "Darwin" ]];then
    go build -v -o fish .
elif [[ "$os" == "Linux" ]];then
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o bin/main .
else
        echo "not support os"
fi
