#!/bin/bash

# lock
touch /tmp/install-java.lock

# install java
apt-get clean && \
	apt-get update --fix-missing && \
	apt-get install -y --fix-missing default-jdk
ln -s /usr/bin/java /usr/local/bin/java

# unlock
rm /tmp/install-java.lock
