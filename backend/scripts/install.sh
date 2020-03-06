#/bin/bash

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
