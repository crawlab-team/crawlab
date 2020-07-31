#!/bin/bash

# fail immediately if error
set -e

# lock global
touch /tmp/install.lock

# lock
touch /tmp/install-java.lock

# install java
apt-get clean
apt-get update --fix-missing
apt-get install -y --fix-missing default-jdk
ln -s /usr/bin/java /usr/local/bin/java

# unlock
rm /tmp/install-java.lock

# unlock global
rm /tmp/install.lock
