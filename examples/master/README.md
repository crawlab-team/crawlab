# Master 节点镜像制作

在Dockerfile里面的二进制包，需要手动在源码目录下进行构建然后再放进来。

## Linux 二进制包构建
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o crawlab main.go
```