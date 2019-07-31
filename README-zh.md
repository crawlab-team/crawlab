# Crawlab

![](http://114.67.75.98:8081/buildStatus/icon?job=crawlab%2Fdevelop)
![](https://img.shields.io/badge/版本-v0.3.0-blue.svg)
<a href="https://github.com/tikazyq/crawlab/blob/master/LICENSE" target="_blank">
    <img src="https://img.shields.io/badge/License-BSD-blue.svg">
</a>

中文 | [English](https://github.com/tikazyq/crawlab)

[安装](#安装) | [运行](#运行) | [截图](#截图) | [架构](#架构) | [集成](#与其他框架的集成) | [比较](#与其他框架比较) | [相关文章](#相关文章) | [社区&赞助](#社区--赞助)

基于Golang的分布式爬虫管理平台，支持Python、NodeJS、Go、Java、PHP等多种编程语言以及多种爬虫框架。

[查看演示 Demo](http://114.67.75.98:8080) | [文档](https://tikazyq.github.io/crawlab-docs)

## 安装

三种方式:
1. [Docker](https://tikazyq.github.io/crawlab/Installation/Docker.md)（推荐）
2. [直接部署](https://tikazyq.github.io/crawlab/Installation/Direct.md)（了解内核）

### 要求（Docker）
- Docker 18.03+
- Redis
- MongoDB 3.6+

### 要求（直接部署）
- Go 1.12+
- Node 8.12+
- Redis
- MongoDB 3.6+

## 运行

### Docker

运行主节点示例。`192.168.99.1`是在Docker Machine网络中的宿主机IP地址。`192.168.99.100`是Docker主节点的IP地址。

```bash
docker run -d --rm --name crawlab \
        -e CRAWLAB_REDIS_ADDRESS=192.168.99.1:6379 \
        -e CRAWLAB_MONGO_HOST=192.168.99.1 \
        -e CRAWLAB_SERVER_MASTER=Y \
        -e CRAWLAB_API_ADDRESS=192.168.99.100:8000 \
        -e CRAWLAB_SPIDER_PATH=/app/spiders \
        -p 8080:8080 \
        -p 8000:8000 \
        -v /var/logs/crawlab:/var/logs/crawlab \
        tikazyq/crawlab:0.3.0
```

当然也可以用`docker-compose`来一键启动，甚至不用配置MongoDB和Redis数据库。在当前目录中创建`docker-compose.yml`文件，输入以下内容。

```bash
version: '3.3'
services:
  master: 
    image: tikazyq/crawlab:latest
    container_name: crawlab-master
    environment:
      CRAWLAB_API_ADDRESS: "192.168.99.100:8000"
      CRAWLAB_SERVER_MASTER: "Y"
      CRAWLAB_MONGO_HOST: "mongo"
      CRAWLAB_REDIS_ADDRESS: "redis:6379"
    ports:    
      - "8080:8080" # frontend
      - "8000:8000" # backend
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

然后执行以下命令，Crawlab主节点＋MongoDB＋Redis就启动了。打开`http://localhost:8080`就能看到界面。

```bash
docker-compose up
```

Docker部署的详情，请见[相关文档](https://tikazyq.github.io/crawlab/Installation/Docker.md)。

### 直接部署

请参考[相关文档](https://tikazyq.github.io/crawlab/Installation/Direct.md)。

## 截图

#### 登录

![](https://crawlab.oss-cn-hangzhou.aliyuncs.com/v0.3.0/login.png?v0.3.0)

#### 首页

![](https://crawlab.oss-cn-hangzhou.aliyuncs.com/v0.3.0/home.png?v0.3.0)

#### 节点列表

![](https://crawlab.oss-cn-hangzhou.aliyuncs.com/v0.3.0/node-list.png?v0.3.0)

#### 节点拓扑图

![](https://crawlab.oss-cn-hangzhou.aliyuncs.com/v0.3.0/node-network.png?v0.3.0)

#### 爬虫列表

![](https://crawlab.oss-cn-hangzhou.aliyuncs.com/v0.3.0/spider-list.png?v0.3.0)

#### 爬虫概览

![](https://crawlab.oss-cn-hangzhou.aliyuncs.com/v0.3.0/spider-overview.png?v0.3.0)

#### 爬虫分析

![](https://crawlab.oss-cn-hangzhou.aliyuncs.com/v0.3.0/spider-analytics.png?v0.3.0)

#### 爬虫文件

![](https://crawlab.oss-cn-hangzhou.aliyuncs.com/v0.3.0/spider-file.png?v0.3.0)

#### 任务详情 - 抓取结果

![](https://crawlab.oss-cn-hangzhou.aliyuncs.com/v0.3.0/task-results.png?v0.3.0_1)

#### 定时任务

![](https://crawlab.oss-cn-hangzhou.aliyuncs.com/v0.3.0/schedule.png?v0.3.0)

## 架构

Crawlab的架构包括了一个主节点（Master Node）和多个工作节点（Worker Node），以及负责通信和数据储存的Redis和MongoDB数据库。

![](https://crawlab.oss-cn-hangzhou.aliyuncs.com/v0.3.0/architecture.png)

前端应用向主节点请求数据，主节点通过MongoDB和Redis来执行任务派发调度以及部署，工作节点收到任务之后，开始执行爬虫任务，并将任务结果储存到MongoDB。架构相对于`v0.3.0`之前的Celery版本有所精简，去除了不必要的节点监控模块Flower，节点监控主要由Redis完成。

### 主节点

主节点是整个Crawlab架构的核心，属于Crawlab的中控系统。

主节点主要负责以下功能:
1. 爬虫任务调度
2. 工作节点管理和通信
3. 爬虫部署
4. 前端以及API服务
5. 执行任务（可以将主节点当成工作节点）

主节点负责与前端应用进行通信，并通过Redis将爬虫任务派发给工作节点。同时，主节点会同步（部署）爬虫给工作节点，通过Redis和MongoDB的GridFS。

### 工作节点

工作节点的主要功能是执行爬虫任务和储存抓取数据与日志，并且通过Redis的`PubSub`跟主节点通信。通过增加工作节点数量，Crawlab可以做到横向扩展，不同的爬虫任务可以分配到不同的节点上执行。

### MongoDB

MongoDB是Crawlab的运行数据库，储存有节点、爬虫、任务、定时任务等数据，另外GridFS文件储存方式是主节点储存爬虫文件并同步到工作节点的中间媒介。

### Redis

Redis是非常受欢迎的Key-Value数据库，在Crawlab中主要实现节点间数据通信的功能。例如，节点会将自己信息通过`HSET`储存在Redis的`nodes`哈希列表中，主节点根据哈希列表来判断在线节点。

### 前端

前端是一个基于[Vue-Element-Admin](https://github.com/PanJiaChen/vue-element-admin)的单页应用。其中重用了很多Element-UI的控件来支持相应的展示。

## 与其他框架的集成

爬虫任务本质上是由一个shell命令来实现的。任务ID将以环境变量`CRAWLAB_TASK_ID`的形式存在于爬虫任务运行的进程中，并以此来关联抓取数据。另外，`CRAWLAB_COLLECTION`是Crawlab传过来的所存放collection的名称。

在爬虫程序中，需要将`CRAWLAB_TASK_ID`的值以`task_id`作为可以存入数据库中`CRAWLAB_COLLECTION`的collection中。这样Crawlab就知道如何将爬虫任务与抓取数据关联起来了。当前，Crawlab只支持MongoDB。

### 集成Scrapy

以下是Crawlab跟Scrapy集成的例子，利用了Crawlab传过来的task_id和collection_name。

```python
import os
from pymongo import MongoClient

MONGO_HOST = '192.168.99.100'
MONGO_PORT = 27017
MONGO_DB = 'crawlab_test'

# scrapy example in the pipeline
class JuejinPipeline(object):
    mongo = MongoClient(host=MONGO_HOST, port=MONGO_PORT)
    db = mongo[MONGO_DB]
    col_name = os.environ.get('CRAWLAB_COLLECTION')
    if not col_name:
        col_name = 'test'
    col = db[col_name]

    def process_item(self, item, spider):
        item['task_id'] = os.environ.get('CRAWLAB_TASK_ID')
        self.col.save(item)
        return item
```

## 与其他框架比较

现在已经有一些爬虫管理框架了，因此为啥还要用Crawlab？

因为很多现有当平台都依赖于Scrapyd，限制了爬虫的编程语言以及框架，爬虫工程师只能用scrapy和python。当然，scrapy是非常优秀的爬虫框架，但是它不能做一切事情。

Crawlab使用起来很方便，也很通用，可以适用于几乎任何主流语言和框架。它还有一个精美的前端界面，让用户可以方便的管理和运行爬虫。

|框架 | 类型 | 分布式 | 前端 | 依赖于Scrapyd |
|:---:|:---:|:---:|:---:|:---:|
| [Crawlab](https://github.com/tikazyq/crawlab) | 管理平台 | Y | Y | N
| [ScrapydWeb](https://github.com/my8100/scrapydweb) | 管理平台 | Y | Y | Y
| [SpiderKeeper](https://github.com/DormyMo/SpiderKeeper) | 管理平台 | Y | Y | Y
| [Gerapy](https://github.com/Gerapy/Gerapy) | 管理平台 | Y | Y | Y
| [Scrapyd](https://github.com/scrapy/scrapyd) | 网络服务 | Y | N | N/A

## 相关文章

- [爬虫管理平台Crawlab部署指南（Docker and more）](https://juejin.im/post/5d01027a518825142939320f)
- [[爬虫手记] 我是如何在3分钟内开发完一个爬虫的](https://juejin.im/post/5ceb4342f265da1bc8540660)
- [手把手教你如何用Crawlab构建技术文章聚合平台(二)](https://juejin.im/post/5c92365d6fb9a070c5510e71)
- [手把手教你如何用Crawlab构建技术文章聚合平台(一)](https://juejin.im/user/5a1ba6def265da430b7af463/posts)

**注意: v0.3.0版本已将基于Celery的Python版本切换为了Golang版本，如何部署请参照文档**

## 社区 & 赞助

如果您觉得Crawlab对您的日常开发或公司有帮助，请加作者微信 tikazyq1 并注明"Crawlab"，作者会将你拉入群。或者，您可以扫下方支付宝二维码给作者打赏去升级团队协作软件或买一杯咖啡。

<p align="center">
    <img src="https://crawlab.oss-cn-hangzhou.aliyuncs.com/gitbook/qrcode.png" height="360">
    <img src="https://crawlab.oss-cn-hangzhou.aliyuncs.com/gitbook/payment.jpg" height="360">
</p>
