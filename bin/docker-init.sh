#!/bin/bash

if [ "${CRAWLAB_NODE_MASTER}" = "Y" ]; then
  # start master
  /bin/bash /app/bin/docker-start-master.sh

  # node type
  echo "node type: master"

  # start crawlab
  crawlab-server master
else
  # node type
  echo "node type: worker"

  # start crawlab
  crawlab-server worker
fi
