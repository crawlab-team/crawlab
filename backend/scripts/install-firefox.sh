#!/bin/bash

# fail immediately if error
set -e

# lock global
touch /tmp/install.lock

# lock
touch /tmp/install-firefox.lock

apt-get update
apt-get -y install  firefox ttf-wqy-microhei ttf-wqy-zenhei xfonts-wqy
apt-get -y install  libcanberra-gtk3-module

# unlock global
rm /tmp/install.lock

# unlock
rm /tmp/install-firefox.lock