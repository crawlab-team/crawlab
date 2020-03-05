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
	/bin/sh /app/backend/scripts/install-nodejs.sh >> /var/log/install-nodejs.sh.log 2>&1 &
fi

# install languages: Java
if [ "${CRAWLAB_SERVER_LANG_JAVA}" = "Y" ];
then
	echo "installing java"
	/bin/sh /app/backend/scripts/install-java.sh >> /var/log/install-java.sh.log 2>&1 &
fi

# generate ssh
ssh-keygen -q -t rsa -N "" -f ${HOME}/.ssh/id_rsa

# ssh config
touch ${HOME}/.ssh/config && chmod 600 ${HOME}/.ssh/config
cat > ${HOME}/.ssh/config <<EOF
Host *
  StrictHostKeyChecking no
EOF

# start backend
crawlab-server