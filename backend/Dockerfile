FROM golang:1.12-alpine AS build

WORKDIR /go/src/app
COPY . .

ENV GO111MODULE on
ENV GOPROXY https://mirrors.aliyun.com/goproxy/

RUN go mod vendor
RUN go install -v ./...

FROM alpine:latest
WORKDIR /root
COPY --from=build /go/src/app .
COPY --from=build /go/bin/crawlab /usr/local/bin

EXPOSE 8000
CMD ["crawlab"]
