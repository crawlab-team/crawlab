docker run --restart always --name crawlab \
        -e CRAWLAB_REDIS_ADDRESS=192.168.99.1:6379 \
        -e CRAWLAB_MONGO_HOST=192.168.99.1 \
        -e CRAWLAB_SERVER_MASTER=N \
        -v /var/logs/crawlab:/var/logs/crawlab \
        tikazyq/crawlab:0.3.0