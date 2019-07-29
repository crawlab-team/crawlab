docker run -d --rm --name crawlab \
        -e CRAWLAB_REDIS_ADDRESS=192.168.99.1:6379 \
        -e CRAWLAB_MONGO_HOST=192.168.99.1 \
        -e CRAWLAB_SERVER_MASTER=Y \
        -e CRAWLAB_API_ADDRESS=192.168.99.100:8000 \
        -e CRAWLAB_SPIDER_PATH=/app/spiders \
        -p 8080:8080 \
        -p 8000:8000 \
        -v /var/logs/crawlab:/var/logs/crawlab \
        tikazyq/crawlab:0.3.0