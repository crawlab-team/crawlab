#!/bin/sh

# replace default api path to new one
if [ "${CRAWLAB_API_ADDRESS}" = "" ]; 
then
	:
else
	jspath=`ls /app/dist/js/app.*.js`
	cp ${jspath} ${jspath}.bak
	sed -i "s/localhost:8000/${CRAWLAB_API_ADDRESS}/g" ${jspath}
fi

# start nginx
service nginx start

crawlab