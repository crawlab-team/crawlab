#!/bin/bash
COMMIT_HASH=$(git rev-parse HEAD)
TIMESTAMP=$(date +%Y%m%d%H%M%S)
echo "v0.0.0-$TIMESTAMP-$COMMIT_HASH"