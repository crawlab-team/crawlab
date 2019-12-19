#!/bin/sh

# replace default api path to new one
if [ "${CRAWLAB_API_ADDRESS}" = "" ]; 
then
	:
else
	jspath=`ls /app/dist/js/app.*.js`
	sed -i "s?http://localhost:8000?${CRAWLAB_API_ADDRESS}?g" ${jspath}
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

# wait for mongo service to be ready
#/app/wait-for-it.sh $CRAWLAB_MONGO_HOST:$CRAWLAB_MONGO_PORT

# wait for redis service to be ready
#/app/wait-for-it.sh $CRAWLAB_REDIS_ADDRESS:$CRAWLAB_REDIS_PORT

crawlab