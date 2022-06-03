FROM golang:1.16

WORKDIR /go/src/app
COPY ./backend .

ENV GO111MODULE on
#ENV GOPROXY https://goproxy.io

RUN go mod tidy \
  && go install -v ./...

