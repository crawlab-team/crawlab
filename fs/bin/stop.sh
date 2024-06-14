#!/bin/sh

kill -9 `ps axu|grep weed|grep -v grep|awk '{print \$2}'|xargs`