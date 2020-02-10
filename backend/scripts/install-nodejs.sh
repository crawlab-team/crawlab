#!/bin/env bash

# install nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.35.2/install.sh | bash
export NVM_DIR="$([ -z "${XDG_CONFIG_HOME-}" ] && printf %s "${HOME}/.nvm" || printf %s "${XDG_CONFIG_HOME}/nvm")"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" # This loads nvm

# install Node.js v8.12
nvm install 8.12

# create soft links
ln -s $HOME/.nvm/versions/node/v8.12.0/bin/npm /usr/local/bin/npm
ln -s $HOME/.nvm/versions/node/v8.12.0/bin/node /usr/local/bin/node

# environments manipulation
export NODE_PATH=$HOME.nvm/versions/node/v8.12.0/lib/node_modules
export PATH=$NODE_PATH:$PATH

# install default dependencies
npm config set PUPPETEER_DOWNLOAD_HOST=https://npm.taobao.org/mirrors
npm install puppeteer-chromium-resolver crawlab-sdk  -g --unsafe-perm=true
