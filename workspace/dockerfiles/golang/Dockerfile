FROM golang:buster
RUN go env -w GOPROXY=https://goproxy.io,https://goproxy.cn && \
    go env -w GO111MODULE="on"
WORKDIR /tools
RUN go get github.com/cosmtrek/air
WORKDIR /backend
RUN rm -rf /tools
# set as non-interactive
ENV DEBIAN_FRONTEND noninteractive

# set CRAWLAB_IS_DOCKER
ENV CRAWLAB_IS_DOCKER Y
# install packages
RUN chmod 777 /tmp \
    && sed -i 's#http://deb.debian.org#https://mirrors.tuna.tsinghua.edu.cn#g' /etc/apt/sources.list \
	&& apt-get update \
	&& apt-get install -y curl net-tools iputils-ping ntp ntpdate python3 python3-pip dumb-init  \
	&& ln -s /usr/bin/pip3 /usr/local/bin/pip \
	&& ln -s /usr/bin/python3 /usr/local/bin/python

RUN pip config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple
# install backend
RUN pip install scrapy pymongo bs4 requests crawlab-sdk scrapy-splash
RUN mkdir /spiders && chmod -R 0755 /spiders
VOLUME /backend
EXPOSE 8080
