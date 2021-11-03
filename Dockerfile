FROM golang:1.15 AS backend-build

WORKDIR /go/src/app
COPY ./backend .

ENV GO111MODULE on
#ENV GOPROXY https://goproxy.io

RUN go mod tidy \
  && go install -v ./...

FROM node:12 AS frontend-build

ADD ./frontend /app
WORKDIR /app
RUN rm /app/.npmrc

# install frontend
RUN yarn install && yarn run build:docker

# images
FROM ubuntu:latest

# set as non-interactive
ENV DEBIAN_FRONTEND noninteractive

# set CRAWLAB_IS_DOCKER
ENV CRAWLAB_IS_DOCKER Y

# install packages
RUN chmod 777 /tmp \
	&& apt-get update \
	&& apt-get install -y curl git net-tools iputils-ping ntp ntpdate python3 python3-pip nginx wget dumb-init cloc \
	&& ln -s /usr/bin/pip3 /usr/local/bin/pip \
	&& ln -s /usr/bin/python3 /usr/local/bin/python

# install seaweedfs
RUN wget https://github.com/chrislusf/seaweedfs/releases/download/2.59/linux_amd64.tar.gz \
  && tar -zxf linux_amd64.tar.gz \
  && cp weed /usr/local/bin

# install backend
RUN pip install scrapy pymongo bs4 requests
RUN pip install crawlab-sdk==0.6.b20210729-1634

# add files
COPY ./backend/conf /app/backend/conf
COPY ./nginx /app/nginx
COPY ./docker_init.sh /app/docker_init.sh

# copy backend files
RUN mkdir -p /opt/bin
COPY --from=backend-build /go/bin/crawlab /opt/bin
RUN cp /opt/bin/crawlab /usr/local/bin/crawlab-server

# copy frontend files
COPY --from=frontend-build /app/dist /app/dist

# copy nginx config files
COPY ./nginx/crawlab.conf /etc/nginx/conf.d

# working directory
WORKDIR /app/backend

# timezone environment
ENV TZ Asia/Shanghai

# language environment
ENV LC_ALL C.UTF-8
ENV LANG C.UTF-8

# frontend port
EXPOSE 8080

# backend port
EXPOSE 8000

# start backend
CMD ["/bin/bash", "/app/docker_init.sh"]
