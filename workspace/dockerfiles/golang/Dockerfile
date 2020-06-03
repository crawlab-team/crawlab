FROM golang:buster
RUN go env -w GOPROXY=https://goproxy.io,https://goproxy.cn && \
    go env -w GO111MODULE="on"
WORKDIR /tools
RUN go get github.com/cosmtrek/air
WORKDIR /backend
RUN rm -rf /tools
VOLUME /backend
EXPOSE 8080