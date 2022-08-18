FROM golang:1.16

RUN go env -w GOPROXY=https://goproxy.io,https://goproxy.cn && \
    go env -w GO111MODULE="on"

WORKDIR /tools
RUN go get github.com/cosmtrek/air

WORKDIR /backend
RUN rm -rf /tools

# set as non-interactive
ENV DEBIAN_FRONTEND noninteractive

# install packages
RUN chmod 777 /tmp \
	&& apt-get update \
	&& apt-get install -y curl git net-tools iputils-ping ntp ntpdate nginx wget dumb-init cloc

# install python
RUN apt-get install -y python3 python3-pip \
	&& ln -s /usr/bin/pip3 /usr/local/bin/pip \
	&& ln -s /usr/bin/python3 /usr/local/bin/python

# install golang
RUN curl -OL https://storage.googleapis.com/golang/go1.16.7.linux-amd64.tar.gz \
	&& tar -C /usr/local -xvf go1.16.7.linux-amd64.tar.gz \
	&& ln -s /usr/local/go/bin/go /usr/local/bin/go \
	&& rm go1.16.7.linux-amd64.tar.gz

# install seaweedfs
RUN wget https://github.com/chrislusf/seaweedfs/releases/download/2.76/linux_amd64.tar.gz \
  && tar -zxf linux_amd64.tar.gz \
  && cp weed /usr/local/bin \
  && rm linux_amd64.tar.gz

# install backend
RUN pip install scrapy pymongo bs4 requests -i https://mirrors.aliyun.com/pypi/simple
RUN pip install crawlab-sdk==0.6.b20211024-1207

VOLUME /backend
EXPOSE 8080
