#!/bin/bash

# lock
touch /tmp/install-java.lock

# install java
rm -r /var/lib/apt/lists/*
apt-get update && apt-get install -y -f default-jdk
ln -s /usr/bin/java /usr/local/bin/java

# unlock
rm /tmp/install-java.lock
