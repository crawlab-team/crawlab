#!/bin/bash

if [ -L $0 ]
then
    BASE_DIR=`dirname $(readlink $0)`
else
    BASE_DIR=`dirname $0`
fi
base_path=$(cd $BASE_DIR/..; pwd)

cd $base_path && \
  rm -rf dist | true

cd $base_path && \
  mkdir -p dist/python | true && \
  mkdir -p dist/js | true && \
  mkdir -p dist/ts | true && \
  mkdir -p dist/java | true && \
  mkdir -p dist/csharp | true && \
  mkdir -p dist/php | true && \
  mkdir -p dist/ruby | true

cd $base_path && \
  protoc -I ./proto \
  --go_out=. \
  --go-grpc_out=. \
  --python_out=dist/python \
  --js_out=dist/js \
  --java_out=dist/java \
  --csharp_out=dist/csharp \
  ./proto/**/*.proto

# python
cd $base_path && \
  python3 -m grpc_tools.protoc -I ./proto \
  --grpc_python_out=dist/python \
  ./proto/**/*.proto
