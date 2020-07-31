#!/bin/bash

# fail immediately if error
set -e

# install node.js
if [ "${CRAWLAB_SERVER_LANG_NODE}" = "Y" ];
then
	echo "installing node.js"
	/bin/sh /app/backend/scripts/install-nodejs.sh
	echo "installed node.js"
fi

# install java
if [ "${CRAWLAB_SERVER_LANG_JAVA}" = "Y" ];
then
	echo "installing java"
	/bin/sh /app/backend/scripts/install-java.sh
	echo "installed java"
fi

# install dotnet
if [ "${CRAWLAB_SERVER_LANG_DOTNET}" = "Y" ];
then
	echo "installing dotnet"
	/bin/sh /app/backend/scripts/install-dotnet.sh
	echo "installed dotnet"
fi

# install php
if [ "${CRAWLAB_SERVER_LANG_PHP}" = "Y" ];
then
	echo "installing php"
	/bin/sh /app/backend/scripts/install-php.sh
	echo "installed php"
fi

# install go
if [ "${CRAWLAB_SERVER_LANG_GO}" = "Y" ];
then
	echo "installing go"
	/bin/sh /app/backend/scripts/install-go.sh
	echo "installed go"
fi
