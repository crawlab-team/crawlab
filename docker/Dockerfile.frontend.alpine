FROM alpine:latest

# 安装nginx
RUN mkdir /run/nginx && apk add nginx

# 拷贝编译文件
COPY dist  /app/dist

# 拷贝nginx代理文件
COPY crawlab.conf /etc/nginx/conf.d

# 拷贝执行脚本
COPY docker_init.sh /app/docker_init.sh

# 定义后端API脚本
ENV CRAWLAB_API_ADDRESS http://localhost:8000

EXPOSE 8080

CMD ["/bin/sh", "/app/docker_init.sh"]