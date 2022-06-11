#!/bin/bash

if [ "${CRAWLAB_NODE_MASTER}" = "Y" ]; then
  # start master
  /bin/bash /app/bin/docker-start-master.sh

  # env
  export IS_MASTER=1

  # start crawlab
  crawlab-server master
else
  # env
  export IS_MASTER=0

  # start crawlab
  crawlab-server worker
fi
