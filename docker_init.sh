#!/bin/bash

# replace absolute api url to relative
jspath=`ls /app/dist/js/index.*.js`
sed -i "s?VUE_APP_API_BASE_URL?\/api?g" ${jspath}

# replace default api path to new one
# if [ "${CRAWLAB_API_ADDRESS}" = "" ]; 
# then
# 	:
# else
# 	jspath=`ls /app/dist/js/app.*.js`
# 	sed -i "s?###CRAWLAB_API_ADDRESS###?${CRAWLAB_API_ADDRESS}?g" ${jspath}
# fi

# replace base url
# if [ "${CRAWLAB_BASE_URL}" = "" ];
# then
# 	:
# else
# 	indexpath=/app/dist/index.html
# 	sed -i "s?/js/?${CRAWLAB_BASE_URL}/js/?g" ${indexpath}
# 	sed -i "s?/css/?${CRAWLAB_BASE_URL}/css/?g" ${indexpath}
# fi

#grant script 
# chmod +x /app/backend/scripts/*.sh

# install languages
# if [ "${CRAWLAB_SERVER_LANG_NODE}" = "Y" ] || [ "${CRAWLAB_SERVER_LANG_JAVA}" = "Y" ] || [ "${CRAWLAB_SERVER_LANG_DOTNET}" = "Y" ] || [ "${CRAWLAB_SERVER_LANG_PHP}" = "Y" ] || [ "${CRAWLAB_SERVER_LANG_GO}" = "Y" ];
# then
# 	echo "installing languages"
# 	echo "you can view log at /var/log/install.sh.log"
# 	/bin/sh /app/backend/scripts/install.sh >> /var/log/install.sh.log 2>&1 &
# fi

# generate ssh
ssh-keygen -q -t rsa -N "" -f ${HOME}/.ssh/id_rsa

# ssh config
touch ${HOME}/.ssh/config && chmod 600 ${HOME}/.ssh/config
cat > ${HOME}/.ssh/config <<EOF
Host *
  StrictHostKeyChecking no
EOF

# start backend
if [ "${CRAWLAB_NODE_MASTER}" = "Y" ]; then
	# start nginx
	service nginx start

	# start seaweedfs server
	seaweedfsDataPath=/data/seaweedfs
	if [ -e ${seaweedfsDataPath} ]; then
		:
	else
		mkdir -p ${seaweedfsDataPath}
	fi
	weed server \
		-dir /data \
		-master.dir ${seaweedfsDataPath} \
		-volume.dir.idx ${seaweedfsDataPath} \
		-ip localhost \
		-volume.port 9999 \
		-filer \
		>> /var/log/weed.log 2>&1 &

	# start crawlab
	crawlab-server master
else
	# start crawlab
	crawlab-server worker
fi
