#!/bin/env bash

# lock
touch /tmp/install-nodejs.lock

# install nvm
BASE_DIR=`dirname $0`
/bin/bash ${BASE_DIR}/install-nvm.sh
export NVM_DIR="$([ -z "${XDG_CONFIG_HOME-}" ] && printf %s "${HOME}/.nvm" || printf %s "${XDG_CONFIG_HOME}/nvm")"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" # This loads nvm

# install Node.js v8.12
export NVM_NODEJS_ORG_MIRROR=http://npm.taobao.org/mirrors/node
nvm install 8.12

# create soft links
ln -s $HOME/.nvm/versions/node/v8.12.0/bin/npm /usr/local/bin/npm
ln -s $HOME/.nvm/versions/node/v8.12.0/bin/node /usr/local/bin/node

# environments manipulation
export NODE_PATH=$HOME.nvm/versions/node/v8.12.0/lib/node_modules
export PATH=$NODE_PATH:$PATH

# install chromium
# See https://crbug.com/795759
apt-get update && apt-get install -yq libgconf-2-4

# Install latest chrome dev package and fonts to support major 
# charsets (Chinese, Japanese, Arabic, Hebrew, Thai and a few others)
# Note: this installs the necessary libs to make the bundled version 
# of Chromium that Puppeteer
# installs, work.
apt-get update && apt-get install -y --no-install-recommends gconf-service libasound2 libatk1.0-0 libatk-bridge2.0-0 libc6 libcairo2 libcups2 libdbus-1-3 libexpat1 libfontconfig1 libgcc1 libgconf-2-4 libgdk-pixbuf2.0-0 libglib2.0-0 libgtk-3-0 libnspr4 libpango-1.0-0 libpangocairo-1.0-0 libstdc++6 libx11-6 libx11-xcb1 libxcb1 libxcomposite1 libxcursor1 libxdamage1 libxext6 libxfixes3 libxi6 libxrandr2 libxrender1 libxss1 libxtst6 ca-certificates fonts-liberation libappindicator1 libnss3 lsb-release xdg-utils wget

# install default dependencies
PUPPETEER_DOWNLOAD_HOST=https://npm.taobao.org/mirrors
npm config set puppeteer_download_host=https://npm.taobao.org/mirrors
npm install puppeteer-chromium-resolver crawlab-sdk -g --unsafe-perm=true --registry=https://registry.npm.taobao.org

# unlock
rm /tmp/install-nodejs.lock
