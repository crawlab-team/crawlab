#!/bin/sh
case $1 in
	master)
		cd /opt/crawlab/frontend \
			&& npm run build:prod \
			&& service nginx start \
			&& mongod --fork --logpath /var/log/mongod.log
		redis-server >> /var/log/redis-server.log 2>&1 &
		python $WORK_DIR/crawlab/flower.py >> /opt/crawlab/flower.log 2>&1 &
		python $WORK_DIR/crawlab/worker.py >> /opt/crawlab/worker.log 2>&1 &
		python $WORK_DIR/crawlab/app.py
			;;
	worker)
		python $WORK_DIR/crawlab/app.py >> /opt/crawlab/app.log 2>&1 &
		python $WORK_DIR/crawlab/worker.py	
		;;
esac