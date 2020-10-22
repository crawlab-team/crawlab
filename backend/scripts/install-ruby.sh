#!/bin/bash

apt-get update
apt-get install -y curl

curl -sSL https://get.rvm.io | bash -s stable
source /etc/profile.d/rvm.sh

echo `rvm list known`
rvm install 2.6.1

echo `ruby -v`
