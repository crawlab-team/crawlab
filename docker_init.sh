cd /opt/crawlab/frontend \
	&& npm run build:prod \
	&& python3 $WORK_DIR/manage.py 