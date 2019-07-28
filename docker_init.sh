#!/bin/sh

# replace default api path to new one
jspath=`ls /app/dist/js/app.*.js`
cat ${jspath} | sed "s/localhost:8000/${CRAWLAB_API_ADDRESS}/g" > ${jspath}

# start nginx
service nginx start

crawlab