#!/bin/sh
case $1 in
	master)
		cd $WORK_DIR/frontend \
			&& npm run build:prod \
			&& service nginx start
		python $WORK_DIR/crawlab/flower.py >> /opt/crawlab/flower.log 2>&1 &
		python $WORK_DIR/crawlab/worker.py >> /opt/crawlab/worker.log 2>&1 &
		cd $WORK_DIR/crawlab \
			&& gunicorn --log-level=DEBUG -b 0.0.0.0 -w 8 $WORK_DIR/crawlab/app:app
			;;
	worker)
		python $WORK_DIR/crawlab/app.py >> /opt/crawlab/app.log 2>&1 &
		python $WORK_DIR/crawlab/worker.py	
		;;
esac