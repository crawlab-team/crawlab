#!/bin/bash
# lock global
touch /tmp/install.lock

# lock
touch /tmp/install-go.lock

# install golang
apt-get update
apt-get install -y golang

# environment variables
export GOPROXY=https://goproxy.cn
export GOPATH=/opt/go

# unlock global
rm /tmp/install.lock

# unlock
rm /tmp/install-go.lock