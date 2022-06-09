#!/bin/bash

if [ "${CRAWLAB_NODE_MASTER}" = "Y" ]; then
  # start master
  /bin/bash /app/bin/docker-start-master.sh

  # start crawlab
  crawlab-server master
else
  # start crawlab
  crawlab-server worker
fi
