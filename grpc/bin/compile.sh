#!/bin/bash

if [ -L $0 ]
then
    BASE_DIR=`dirname $(readlink $0)`
else
    BASE_DIR=`dirname $0`
fi
base_path=$(cd $BASE_DIR/..; pwd)

cd $base_path && \
  protoc -I ./proto \
  --go_out=. \
  --go-grpc_out=. \
  ./proto/**/*.proto