# 本地开发环境worker节点制作
由于master和worker节点的存储信息是在redis上，并且使用节点所在的mac地址作为key，所以在开发本地需要启动master和worker节点会比较麻烦。
这里是一个运行worker节点的一个例子。

基本思路是worker节点所需的依赖制作成一个镜像，然后把crawlab编译成二进制包，接着把配置文件和二进制包通过volumes的形式挂载到容器内部。
这样就可以正常的运行worker节点了。之后对于容器编排的worker节点，可以直接把该镜像当成worker节点的基础镜像。

### 制作二进制包
在`backend`目录下执行以下命令，生成二进制包
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o crawlab main.go
```


### 构建worker镜像
```
docker build -t crawlab:worker .
```

### 运行worker节点
```
docker-compose up -d
```

如果在多台服务器使用`docker-compose.yml`进行编排，可能出现节点注册不上的问题，因为mac地址冲突了。
可以使用`networks`定义当前节点的IP段，这样就可以正常注册到redis