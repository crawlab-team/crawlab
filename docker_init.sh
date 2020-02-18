#!/bin/bash

# replace default api path to new one
if [ "${CRAWLAB_API_ADDRESS}" = "" ]; 
then
	:
else
	jspath=`ls /app/dist/js/app.*.js`
	sed -i "s?###CRAWLAB_API_ADDRESS###?${CRAWLAB_API_ADDRESS}?g" ${jspath}
fi

# replace base url
if [ "${CRAWLAB_BASE_URL}" = "" ];
then
	:
else
	indexpath=/app/dist/index.html
	sed -i "s?/js/?${CRAWLAB_BASE_URL}/js/?g" ${indexpath}
	sed -i "s?/css/?${CRAWLAB_BASE_URL}/css/?g" ${indexpath}
fi

# start nginx
service nginx start

# install languages: Node.js
if [ "${CRAWLAB_SERVER_LANG_NODE}" = "Y" ];
then
	echo "installing node.js"
	/bin/sh /app/backend/scripts/install-nodejs.sh
fi

# generate ssh
ssh-keygen -q -t rsa -N "" -f ${HOME}/.ssh/id_rsa

# start backend
crawlab