#!/bin/env bash

# lock
touch /tmp/install-java.lock

# install java
apt-get update
apt-get install -y default-jdk
ln -s /usr/bin/java /usr/local/bin/java

# unlock
rm /tmp/install-java.lock
