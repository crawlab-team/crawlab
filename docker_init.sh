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

#grant script 
chmod +x /app/backend/scripts/*.sh

# install languages
if [ "${CRAWLAB_SERVER_LANG_NODE}" = "Y" ] || [ "${CRAWLAB_SERVER_LANG_JAVA}" = "Y" ] || [ "${CRAWLAB_SERVER_LANG_DOTNET}" = "Y" ] || [ "${CRAWLAB_SERVER_LANG_PHP}" = "Y" ] || [ "${CRAWLAB_SERVER_LANG_GO}" = "Y" ];
then
	echo "installing languages"
	echo "you can view log at /var/log/install.sh.log"
	/bin/sh /app/backend/scripts/install.sh >> /var/log/install.sh.log 2>&1 &
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