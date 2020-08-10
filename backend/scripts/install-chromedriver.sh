#!/bin/bash

# fail immediately if error
set -e

# lock global
touch /tmp/install.lock

# lock
touch /tmp/install-chromedriver.lock

export DEBIAN_FRONTEND=noninteractive
apt-get update
apt-get install unzip
DL=https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
curl -sL "$DL" > /tmp/chrome.deb
apt install --no-install-recommends --no-install-suggests -y /tmp/chrome.deb
CHROMIUM_FLAGS='--no-sandbox --disable-dev-shm-usage'
sed -i '${s/$/'" $CHROMIUM_FLAGS"'/}' /opt/google/chrome/google-chrome
BASE_URL=https://chromedriver.storage.googleapis.com
VERSION=$(curl -sL "$BASE_URL/LATEST_RELEASE")
curl -sL "$BASE_URL/$VERSION/chromedriver_linux64.zip" -o /tmp/driver.zip
unzip /tmp/driver.zip
chmod 755 chromedriver
mv chromedriver /usr/local/bin

# unlock global
rm /tmp/install.lock

# unlock
rm /tmp/install-chromedriver.lock