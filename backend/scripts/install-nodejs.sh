#!/bin/env bash

# install nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.35.2/install.sh | bash
export NVM_DIR="$([ -z "${XDG_CONFIG_HOME-}" ] && printf %s "${HOME}/.nvm" || printf %s "${XDG_CONFIG_HOME}/nvm")"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" # This loads nvm

# install Node.js v8.12
# export NVM_NODEJS_ORG_MIRROR=http://npm.taobao.org/mirrors/node
nvm install 8.12

# create soft links
ln -s $HOME/.nvm/versions/node/v8.12.0/bin/npm /usr/local/bin/npm
ln -s $HOME/.nvm/versions/node/v8.12.0/bin/node /usr/local/bin/node

# environments manipulation
export NODE_PATH=$HOME.nvm/versions/node/v8.12.0/lib/node_modules
export PATH=$NODE_PATH:$PATH

# install apt dependencies
apt-get install -y wget
wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - \
    && sh -c 'echo "deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list' \
    && apt-get update \
    && apt-get install -y google-chrome-unstable \
      --no-install-recommends \
    && rm -rf /var/lib/apt/lists/*

# install default dependencies
export PUPPETEER_DOWNLOAD_HOST=https://npm.taobao.org/mirrors
npm install puppeteer-chromium-resolver crawlab-sdk -g --unsafe-perm=true
