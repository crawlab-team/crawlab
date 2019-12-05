## Docker安装部署

这应该是部署应用的最方便也是最节省时间的方式了。在最近的一次版本更新[v0.3.0](https://github.com/tikazyq/crawlab/releases/tag/v0.3.0)中，我们发布了Golang版本，并且支持Docker部署。下面将一步一步介绍如何使用Docker来部署Crawlab。

对Docker不了解的开发者，可以参考一下这篇文章（[9102 年了，学点 Docker 知识](https://juejin.im/post/5c2c69cee51d450d9707236e)）做进一步了解。简单来说，Docker可以利用已存在的镜像帮助构建一些常用的服务和应用，例如Nginx、MongoDB、Redis等等。用Docker运行一个MongoDB服务仅需`docker run -d --name mongo -p 27017:27017 mongo`一行命令。如何安装Docker跟操作系统有关，这里就不展开讲了，需要的同学自行百度一下相关教程。

### 下载镜像

我们已经在[DockerHub](https://hub.docker.com/r/tikazyq/crawlab)上构建了Crawlab的镜像，开发者只需要将其pull下来使用。在pull 镜像之前，我们需要配置一下镜像源。因为我们在墙内，使用原有的镜像源速度非常感人，因此将使用DockerHub在国内的加速器。创建`/etc/docker/daemon.json`文件，在其中输入如下内容。

```json
{
  "registry-mirrors": ["https://registry.docker-cn.com"]
}
```

这样的话，pull镜像的速度会比不改变镜像源的速度快很多。

执行以下命令将Crawlab的镜像下载下来。镜像大小大概在几百兆，因此下载需要几分钟时间。

```bash
docker pull tikazyq/crawlab:latest
```

### 运行Docker容器

之前的版本需要更改配置文件来配置数据库等参数，非常麻烦。在最近的版本`v0.3.0`中，我们实现了用环境变量来替代配置文件，简化了配置步骤。

运行以下命令启动主节点。

```bash
docker run -d --restart always --name crawlab \
        -e CRAWLAB_REDIS_ADDRESS=192.168.99.1 \
        -e CRAWLAB_MONGO_HOST=192.168.99.1 \
        -e CRAWLAB_SERVER_MASTER=Y \
        -e CRAWLAB_API_ADDRESS=192.168.99.100:8000 \
        -p 8080:8080 \
        -p 8000:8000 \
        -v /var/logs/crawlab:/var/logs/crawlab \
        tikazyq/crawlab:0.3.0
```

其中，我们设置了Redis和MongoDB的地址，分别通过`CRAWLAB_REDIS_ADDRESS`和`CRAWLAB_MONGO_HOST`参数。`CRAWLAB_SERVER_MASTER`设置为`Y`表示启动的是主节点（该参数默认是为`N`，表示为工作节点）。`CRAWLAB_API_ADDRESS`是前端的API地址，请将这个设置为公网能访问到主节点的地址，`8000`是API端口。环境变量配置详情请见[配置Crawlab](../Config/README.md)，您可以根据自己的要求来进行配置。

此外，我们通过`-p`参数映射了8080端口（Nginx前端静态文件）以及8000端口（后端API）到宿主机。我们将任务日志目录`/var/logs/crawlab`映射出来，保证Docker重启时不会丢失日志文件（当然，我们现在用文件系统来保存日志的方式可能不是一个很好的解决方案，如果您有更好的建议，请在Github上提Issue或者Pull Request）。

您可能好奇为什么我们用`192.168.99.1`，而不是`localhost`。这是因为我这里的例子是用的Docker Machine，它会在宿主机创建一个`192.168.99.*`的网络，而`192.168.99.1`是宿主机的IP地址，`192.168.99.100`就是该Docker Container的地址。因此，这里的启动配置表示，我们启动的主节点连接的是宿主机的Redis和MongoDB，而API地址为该主节点地址。当然，为了方便配置，我们可以用`docker-compose`来管理Docker集群，让他们在同一个网络中，后面我们会介绍。

类似，我们也可以启动工作节点。

```bash
docker run --restart always --name crawlab \
        -e CRAWLAB_REDIS_ADDRESS=192.168.99.1 \
        -e CRAWLAB_MONGO_HOST=192.168.99.1 \
        -e CRAWLAB_SERVER_MASTER=N \
        -v /var/logs/crawlab:/var/logs/crawlab \
        tikazyq/crawlab:latest
```

这里，我们将`CRAWLAB_SERVER_MASTER`设置为`N`，表示它为工作节点（切勿设置多个节点为`Y`，这可能会导致无法预测的问题）。

以上两个Docker启动的命令在Github上，详情请见[Examples](https://github.com/tikazyq/crawlab/tree/master/examples)。

### Docker-Compose

当然，也可以用`docker-compose`的方式来部署。`docker-compose`是一个集群管理方式，可以利用名为`docker-compose.yml`的`yaml`文件来定义需要启动的容器，可以是单个，也可以（通常）是多个的。

Crawlab的`docker-compose.yml`定义如下。

```yaml
version: '3.3'
services:
  master: 
    image: tikazyq/crawlab:latest
    container_name: crawlab-master
    environment:
      CRAWLAB_API_ADDRESS: "192.168.99.100:8000"
      CRAWLAB_SERVER_MASTER: "Y"
      CRAWLAB_MONGO_HOST: "mongo"
      CRAWLAB_REDIS_ADDRESS: "redis"
    ports:    
      - "8080:8080" # frontend
      - "8000:8000" # backend
    depends_on:
      - mongo
      - redis
  worker:
    image: tikazyq/crawlab:latest
    container_name: crawlab-worker
    environment:
      CRAWLAB_SERVER_MASTER: "N"
      CRAWLAB_MONGO_HOST: "mongo"
      CRAWLAB_REDIS_ADDRESS: "redis"
    depends_on:
      - mongo
      - redis
  mongo:
    image: mongo:latest
    restart: always
    ports:
      - "27017:27017"
  redis:
    image: redis:latest
    restart: always
    ports:
      - "6379:6379"
```

这里先定义了`master`节点和`worker`节点，也就是Crawlab的主节点和工作节点。`master`和`worker`依赖于`mongo`和`redis`容器，因此在启动之前会同时启动`mongo`和`redis`容器。这样就不需要单独配置`mongo`和`redis`服务了，大大节省了环境配置的时间。

安装`docker-compose`也很简单，大家去网上百度一下有中文教程。英语水平还可以的可以参考一下[官方文档](https://docs.docker.com/compose/)。

安装完`docker-compose`和定义好`docker-compose.yml`后，只需要运行以下命令就可以启动Crawlab。

```bash
docker-compose up
```

同样，在浏览器中输入`http://localhost:8080`就可以看到界面（Docker Machine是`192.168.99.100`）。
