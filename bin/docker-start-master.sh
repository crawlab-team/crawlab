#!/bin/bash

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
