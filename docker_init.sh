#!/bin/sh
cd /opt/crawlab/frontend \
	&& npm run build:prod \
	&& mongod --fork --logpath /var/log/mongod.log \
	&& service nginx start \
	&& service redis-server start \
	&& python3 $WORK_DIR/manage.py $*