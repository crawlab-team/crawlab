#!/bin/bash

# replace default api path to new one
python /app/bin/update_docker_js_api_address.py

crawlab-server server
