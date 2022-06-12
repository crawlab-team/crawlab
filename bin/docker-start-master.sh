#!/bin/bash
# NOTE: deprecated

# replace env
indexPath=/app/dist/index.html
if test -z "$CRAWLAB_BASE_URL"; then
  :
else
  sed -i "s?/js/?${CRAWLAB_BASE_URL}/js/?g" ${indexPath}
  sed -i "s?/css/?${CRAWLAB_BASE_URL}/css/?g" ${indexPath}
  sed -i "s/  <link rel=\"icon\" type=\"image\/x-icon\" href=\"/  <link rel=\"icon\" type=\"image\/x-icon\" href=\"\/${CRAWLAB_BASE_URL}/g" ${indexPath}
  sed -i "s/  <link rel=\"stylesheet\" href=\"/  <link rel=\"stylesheet\" href=\"${CRAWLAB_BASE_URL}\//g" ${indexPath}
  sed -i "s/  window.VUE_APP_API_BASE_URL = '/  window.VUE_APP_API_BASE_URL = '\/${CRAWLAB_BASE_URL}/g" ${indexPath}
fi
if test -z "$CRAWLAB_INIT_BAIDU_TONGJI"; then
  :
else
  sed -i "s/  window.VUE_APP_INIT_BAIDU_TONGJI = ''/  window.VUE_APP_INIT_BAIDU_TONGJI = '${CRAWLAB_INIT_BAIDU_TONGJI}'/g" ${indexPath}
fi
if test -z "$CRAWLAB_INIT_UMENG"; then
  :
else
  sed -i "s/  window.VUE_APP_INIT_UMENG = ''/  window.VUE_APP_INIT_UMENG = '${CRAWLAB_INIT_UMENG}'/g" ${indexPath}
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
  >>/var/log/weed.log 2>&1 &
