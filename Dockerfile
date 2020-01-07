FROM golang:1.12 AS backend-build

WORKDIR /go/src/app
COPY ./backend .

ENV GO111MODULE on
ENV GOPROXY https://goproxy.io

RUN go install -v ./...

FROM node:8.16.0-alpine AS frontend-build

ADD ./frontend /app
WORKDIR /app

# install frontend
RUN npm config set unsafe-perm true
RUN npm install -g yarn && yarn install

RUN npm run build:prod

# images
FROM ubuntu:latest

ADD . /app

# set as non-interactive
ENV DEBIAN_FRONTEND noninteractive

# set CRAWLAB_IS_DOCKER
ENV CRAWLAB_IS_DOCKER Y

# install packages
RUN apt-get update \
	&& apt-get install -y curl git net-tools iputils-ping ntp ntpdate python3 python3-pip \
	&& ln -s /usr/bin/pip3 /usr/local/bin/pip \
	&& ln -s /usr/bin/python3 /usr/local/bin/python

# install backend
RUN pip install scrapy pymongo bs4 requests

# copy backend files
COPY --from=backend-build /go/bin/crawlab /usr/local/bin

# install nginx
RUN apt-get -y install nginx

# copy frontend files
COPY --from=frontend-build /app/dist /app/dist
COPY --from=frontend-build /app/conf/crawlab.conf /etc/nginx/conf.d

# working directory
WORKDIR /app/backend

# frontend port
EXPOSE 8080

# backend port
EXPOSE 8000

# start backend
CMD ["/bin/bash", "/app/docker_init.sh"]
