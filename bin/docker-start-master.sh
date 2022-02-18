# replace base url
if [ "${CRAWLAB_BASE_URL}" = "" ];
then
	:
else
	indexpath=/app/dist/index.html
	sed -i "s?/js/?${CRAWLAB_BASE_URL}/js/?g" ${indexpath}
	sed -i "s?/css/?${CRAWLAB_BASE_URL}/css/?g" ${indexpath}

	sed -i "s/  <link rel=\"icon\" type=\"image\/x-icon\" href=\"/  <link rel=\"icon\" type=\"image\/x-icon\" href=\"\/${CRAWLAB_BASE_URL}/g"  ${indexpath}
	sed -i "s/  <link rel=\"stylesheet\" href=\"/  <link rel=\"stylesheet\" href=\"${CRAWLAB_BASE_URL}\//g"  ${indexpath}
	sed -i "s/  window.VUE_APP_API_BASE_URL = '/  window.VUE_APP_API_BASE_URL = '\/${CRAWLAB_BASE_URL}/g" ${indexpath}
fi

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

