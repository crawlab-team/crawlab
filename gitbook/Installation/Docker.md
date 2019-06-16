## Docker安装部署

这应该是部署应用的最方便也是最节省时间的方式了。在最近的一次版本更新[v0.2.3](https://github.com/tikazyq/crawlab/releases/tag/v0.2.3)中，我们发布了Docker功能，让大家可以利用Docker来轻松部署Crawlab。下面将一步一步介绍如何使用Docker来部署Crawlab。

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

### 更改配置文件

拷贝一份后端配置文件`./crawlab/config/config.py`以及前端配置文件`./frontend/.env.production`到某一个地方。例如我的例子，分别为`/home/yeqing/config.py`和`/home/yeqing/.env.production`。

更改后端配置文件`config.py`，将MongoDB、Redis的指向IP更改为自己数据的值。注意，容器中对应的宿主机的IP地址不是`localhost`，而是`172.17.0.1`（当然也可以用network来做，只是稍微麻烦一些）。更改前端配置文件`.env.production`，将API地址`VUE_APP_BASE_URL`更改为宿主机所在的IP地址，例如`http://192.168.0.8:8000`，这将是前端调用API会用到的URL。

### 运行Docker容器

更改好配置文件之后，接下来就是运行容器了。执行以下命令来启动容器。

```bash
docker run -d --rm --name crawlab \
	-p 8080:8080 \
	-p 8000:8000 \
	-v /home/yeqing/.env.production:/opt/crawlab/frontend/.env.production \
	-v /home/yeqing/config.py:/opt/crawlab/crawlab/config/config.py \
	tikazyq/crawlab master
```

其中，我们映射了8080端口（Nginx前端静态文件）以及8000端口（后端API）到宿主机。另外还将前端配置文件`/home/yeqing/.env.production`和后端配置文件`/home/yeqing/config.py`映射到了容器相应的目录下。传入参数`master`是代表该启动方式为主机启动模式，也就是所有服务（前端、Api、Flower、Worker）都会启动。另外一个模式是`worker`模式，只会启动必要的Api和Worker服务，这个对于分布式部署比较有用。等待大约20-30秒的时间来build前端静态文件，之后就可以打开Crawlab界面地址地址看到界面了。界面地址默认为`http://localhost:8080`。

![](https://user-gold-cdn.xitu.io/2019/6/12/16b4c3ed5dcd6cfc?w=2532&h=1300&f=png&s=146531)

### Docker-Compose

当然，也可以用`docker-compose`的方式来部署。`docker-compose`是一个集群管理方式，可以利用名为`docker-compose.yml`的`yaml`文件来定义需要启动的容器，可以是单个，也可以（通常）是多个的。Crawlab的`docker-compose.yml`定义如下。

```yaml
version: '3.3'
services:
  master: 
    image: tikazyq/crawlab:latest
    container_name: crawlab
    volumns:
      - /home/yeqing/config.py:/opt/crawlab/crawlab/config/config.py # 后端配置文件
      - /home/yeqing/.env.production:/opt/crawlab/frontend/.env.production # 前端配置文件
    ports:    
      - "8080:8080" # nginx
      - "8000:8000" # app
    depends_on:
      - mongo
      - redis
    entrypoint:
      - /bin/sh
      - /opt/crawlab/docker_init.sh
      - master
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

这里先定义了`master`节点，也就是Crawlab的主节点。`master`依赖于`mongo`和`redis`容器，因此在启动之前会同时启动`mongo`和`redis`容器。这样就不需要单独配置`mongo`和`redis`服务了，大大节省了环境配置的时间。

安装`docker-compose`也很简单，大家去网上百度一下就可以了。

安装完`docker-compose`和定义好`docker-compose.yml`后，只需要运行以下命令就可以启动Crawlab。

```bash
docker-compose up
```

同样，在浏览器中输入`http://localhost:8080`就可以看到界面。

### 多节点模式

`docker-compose`的方式很适合多节点部署，在原有的`master`基础上增加几个`worker`节点，达到多节点部署的目的。将`docker-compose.yml`更改为如下内容。

```yaml
version: '3.3'
services:
  master: 
    image: tikazyq/crawlab:latest
    container_name: crawlab
    volumns:
      - /home/yeqing/config.master.py:/opt/crawlab/crawlab/config/config.py # 后端配置文件
      - /home/yeqing/.env.production.master:/opt/crawlab/frontend/.env.production # 前端配置文件
    ports:    
      - "8080:8080" # nginx
      - "8000:8000" # app
    depends_on:
      - mongo
      - redis
    entrypoint:
      - /bin/sh
      - /opt/crawlab/docker_init.sh
      - master
  worker1: 
    image: tikazyq/crawlab:latest
    volumns:
      - /home/yeqing/config.worker.py:/opt/crawlab/crawlab/config/config.py # 后端配置文件
      - /home/yeqing/.env.production.worker:/opt/crawlab/frontend/.env.production # 前端配置文件
    ports:
      - "8001:8000" # app
    depends_on:
      - mongo
      - redis
    entrypoint:
      - /bin/sh
      - /opt/crawlab/docker_init.sh
      - worker
  worker2: 
    image: tikazyq/crawlab:latest
    volumns:
      - /home/yeqing/config.worker.py:/opt/crawlab/crawlab/config/config.py # 后端配置文件
      - /home/yeqing/.env.production.worker:/opt/crawlab/frontend/.env.production # 前端配置文件
    ports:
      - "8002:8000" # app
    depends_on:
      - mongo
      - redis
    entrypoint:
      - /bin/sh
      - /opt/crawlab/docker_init.sh
      - worker
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

这里启动了多增加了两个`worker`节点，以`worker`模式启动。这样，多节点部署，也就是分布式部署就完成了。