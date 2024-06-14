#!/bin/sh
if [ -e ./tmp ]; then
  :
else
  mkdir ./tmp
fi

if [ -x /usr/local/bin/weed ]; then
  weed server \
    -dir ./tmp \
    -master.dir ./tmp \
    -volume.dir.idx ./tmp \
    -ip localhost \
    -ip.bind 0.0.0.0 \
    -filer
else
  ./seaweedfs/weed server \
    -dir ./tmp \
    -master.dir ./tmp \
    -volume.dir.idx ./tmp \
    -ip localhost \
    -ip.bind 0.0.0.0 \
    -filer
fi
