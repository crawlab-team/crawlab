FROM alpine:latest

# 配置工作目录
WORKDIR /opt/crawlab

# 拷贝配置文件
COPY config.yml /opt/crawlab/conf/config.yml

# 拷贝可执行文件
COPY crawlab /usr/local/bin

# 创建spiders文件用于存放爬虫， 授权可执行文件
RUN mkdir -p /opt/crawlab/spiders && chmod +x /usr/local/bin/crawlab

# 指定为Master节点
ENV CRAWLAB_SERVER_MASTER Y

EXPOSE 8000

CMD ["crawlab"]