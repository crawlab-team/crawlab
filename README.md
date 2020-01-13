# Crawlab

<p>
  <a href="https://hub.docker.com/r/tikazyq/crawlab/builds" target="_blank">
    <img src="https://img.shields.io/docker/cloud/build/tikazyq/crawlab.svg?label=build&logo=docker">
  </a>
  <a href="https://hub.docker.com/r/tikazyq/crawlab" target="_blank">
    <img src="https://img.shields.io/docker/pulls/tikazyq/crawlab?label=pulls&logo=docker">
  </a>
  <a href="https://github.com/crawlab-team/crawlab/releases" target="_blank">
    <img src="https://img.shields.io/github/release/crawlab-team/crawlab.svg?logo=github">
  </a>
  <a href="https://github.com/crawlab-team/crawlab/commits/master" target="_blank">
    <img src="https://img.shields.io/github/last-commit/crawlab-team/crawlab.svg">
  </a>
  <a href="https://github.com/crawlab-team/crawlab/issues?q=is%3Aissue+is%3Aopen+label%3Abug" target="_blank">
    <img src="https://img.shields.io/github/issues/crawlab-team/crawlab/bug.svg?label=bugs&color=red">
  </a>
  <a href="https://github.com/crawlab-team/crawlab/issues?q=is%3Aissue+is%3Aopen+label%3Aenhancement" target="_blank">
    <img src="https://img.shields.io/github/issues/crawlab-team/crawlab/enhancement.svg?label=enhancements&color=cyan">
  </a>
  <a href="https://github.com/crawlab-team/crawlab/blob/master/LICENSE" target="_blank">
    <img src="https://img.shields.io/github/license/crawlab-team/crawlab.svg">
  </a>
</p>

[中文](https://github.com/crawlab-team/crawlab/blob/master/README-zh.md) | English

[Installation](#installation) | [Run](#run) | [Screenshot](#screenshot) | [Architecture](#architecture) | [Integration](#integration-with-other-frameworks) | [Compare](#comparison-with-other-frameworks) | [Community & Sponsorship](#community--sponsorship) | [CHANGELOG](https://github.com/crawlab-team/crawlab/blob/master/CHANGELOG.md) | [Disclaimer](https://github.com/crawlab-team/crawlab/blob/master/DISCLAIMER.md)

Golang-based distributed web crawler management platform, supporting various languages including Python, NodeJS, Go, Java, PHP and various web crawler frameworks including Scrapy, Puppeteer, Selenium.

[Demo](http://crawlab.cn/demo) | [Documentation](http://docs.crawlab.cn)

## Installation

Two methods:
1. [Docker](http://docs.crawlab.cn/Installation/Docker.html) (Recommended)
2. [Direct Deploy](http://docs.crawlab.cn/Installation/Direct.html) (Check Internal Kernel)
3. [Kubernetes](https://juejin.im/post/5e0a02d851882549884c27ad) (Multi-Node Deployment)

### Pre-requisite (Docker)
- Docker 18.03+
- Redis 5.x+
- MongoDB 3.6+
- Docker Compose 1.24+ (optional but recommended)

### Pre-requisite (Direct Deploy)
- Go 1.12+
- Node 8.12+
- Redis 5.x+
- MongoDB 3.6+

## Quick Start

Please open the command line prompt and execute the command beloe. Make sure you have installed `docker-compose` in advance.

```bash
git clone https://github.com/crawlab-team/crawlab
cd crawlab
docker-compose up -d
```

Next, you can look into the `docker-compose.yml` (with detailed config params) and the [Documentation (Chinese)](http://docs.crawlab.cn) for further information. 

## Run

### Docker

Please use `docker-compose` to one-click to start up. By doing so, you don't even have to configure MongoDB and Redis databases. Create a file named `docker-compose.yml` and input the code below.


```yaml
version: '3.3'
services:
  master: 
    image: tikazyq/crawlab:latest
    container_name: master
    environment:
      CRAWLAB_SERVER_MASTER: "Y"
      CRAWLAB_MONGO_HOST: "mongo"
      CRAWLAB_REDIS_ADDRESS: "redis"
    ports:    
      - "8080:8080"
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

Then execute the command below, and Crawlab Master Node + MongoDB + Redis will start up. Open the browser and enter `http://localhost:8080` to see the UI interface.

```bash
docker-compose up
```

For Docker Deployment details, please refer to [relevant documentation](https://tikazyq.github.io/crawlab-docs/Installation/Docker.html).


## Screenshot

#### Login

![](https://raw.githubusercontent.com/tikazyq/crawlab-docs/master/images/login.png)

#### Home Page

![](https://raw.githubusercontent.com/tikazyq/crawlab-docs/master/images/home.png)

#### Node List

![](https://raw.githubusercontent.com/tikazyq/crawlab-docs/master/images/node-list.png)

#### Node Network

![](https://raw.githubusercontent.com/tikazyq/crawlab-docs/master/images/node-network.png)

#### Spider List

![](https://raw.githubusercontent.com/tikazyq/crawlab-docs/master/images/spider-list.png)

#### Spider Overview

![](https://raw.githubusercontent.com/tikazyq/crawlab-docs/master/images/spider-overview.png)

#### Spider Analytics

![](https://raw.githubusercontent.com/tikazyq/crawlab-docs/master/images/spider-analytics.png)

#### Spider Files

![](https://raw.githubusercontent.com/tikazyq/crawlab-docs/master/images/spider-file.png)

#### Task Results

![](https://raw.githubusercontent.com/tikazyq/crawlab-docs/master/images/task-results.png)

#### Cron Job

![](https://raw.githubusercontent.com/tikazyq/crawlab-docs/master/images/schedule.png)

#### Dependency Installation

![](http://static-docs.crawlab.cn/node-install-dependencies.png)

## Architecture

The architecture of Crawlab is consisted of the Master Node and multiple Worker Nodes, and Redis and MongoDB databases which are mainly for nodes communication and data storage.

![](http://static-docs.crawlab.cn/architecture.png)

The frontend app makes requests to the Master Node, which assigns tasks and deploys spiders through MongoDB and Redis. When a Worker Node receives a task, it begins to execute the crawling task, and stores the results to MongoDB. The architecture is much more concise compared with versions before `v0.3.0`. It has removed unnecessary Flower module which offers node monitoring services. They are now done by Redis.

### Master Node

The Master Node is the core of the Crawlab architecture. It is the center control system of Crawlab.

The Master Node offers below services:
1. Crawling Task Coordination;
2. Worker Node Management and Communication;
3. Spider Deployment;
4. Frontend and API Services;
5. Task Execution (one can regard the Master Node as a Worker Node)

The Master Node communicates with the frontend app, and send crawling tasks to Worker Nodes. In the mean time, the Master Node synchronizes (deploys) spiders to Worker Nodes, via Redis and MongoDB GridFS.

### Worker Node

The main functionality of the Worker Nodes is to execute crawling tasks and store results and logs, and communicate with the Master Node through Redis `PubSub`. By increasing the number of Worker Nodes, Crawlab can scale horizontally, and different crawling tasks can be assigned to different nodes to execute.

### MongoDB

MongoDB is the operational database of Crawlab. It stores data of nodes, spiders, tasks, schedules, etc. The MongoDB GridFS file system is the medium for the Master Node to store spider files and synchronize to the Worker Nodes.

### Redis

Redis is a very popular Key-Value database. It offers node communication services in Crawlab. For example, nodes will execute `HSET` to set their info into a hash list named `nodes` in Redis, and the Master Node will identify online nodes according to the hash list.

### Frontend

Frontend is a SPA based on 
[Vue-Element-Admin](https://github.com/PanJiaChen/vue-element-admin). It has re-used many Element-UI components to support corresponding display. 

## Integration with Other Frameworks

A crawling task is actually executed through a shell command. The Task ID will be passed to the crawling task process in the form of environment variable named `CRAWLAB_TASK_ID`. By doing so, the data can be related to a task. Also, another environment variable `CRAWLAB_COLLECTION` is passed by Crawlab as the name of the collection to store results data.

### Scrapy

Below is an example to integrate Crawlab with Scrapy in pipelines. 

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

## Comparison with Other Frameworks

There are existing spider management frameworks. So why use Crawlab? 

The reason is that most of the existing platforms are depending on Scrapyd, which limits the choice only within python and scrapy. Surely scrapy is a great web crawl framework, but it cannot do everything. 

Crawlab is easy to use, general enough to adapt spiders in any language and any framework. It has also a beautiful frontend interface for users to manage spiders much more easily. 

|Framework | Type | Distributed | Frontend | Scrapyd-Dependent |
|:---:|:---:|:---:|:---:|:---:|
| [Crawlab](https://github.com/crawlab-team/crawlab) | Admin Platform | Y | Y | N
| [ScrapydWeb](https://github.com/my8100/scrapydweb) | Admin Platform | Y | Y | Y
| [SpiderKeeper](https://github.com/DormyMo/SpiderKeeper) | Admin Platform | Y | Y | Y
| [Gerapy](https://github.com/Gerapy/Gerapy) | Admin Platform | Y | Y | Y
| [Scrapyd](https://github.com/scrapy/scrapyd) | Web Service | Y | N | N/A

## Contributors
<a href="https://github.com/tikazyq">
  <img src="https://avatars3.githubusercontent.com/u/3393101?s=460&v=4" height="80">
</a>
<a href="https://github.com/wo10378931">
  <img src="https://avatars2.githubusercontent.com/u/8297691?s=460&v=4" height="80">
</a>
<a href="https://github.com/yaziming">
  <img src="https://avatars2.githubusercontent.com/u/54052849?s=460&v=4" height="80">
</a>
<a href="https://github.com/hantmac">
  <img src="https://avatars2.githubusercontent.com/u/7600925?s=460&v=4" height="80">
</a>
<a href="https://github.com/duanbin0414">
  <img src="https://avatars3.githubusercontent.com/u/50389867?s=460&v=4" height="80">
</a>

## Community & Sponsorship

If you feel Crawlab could benefit your daily work or your company, please add the author's Wechat account noting "Crawlab" to enter the discussion group. Or you scan the Alipay QR code below to give us a reward to upgrade our teamwork software or buy a coffee.

<p align="center">
    <img src="https://crawlab.oss-cn-hangzhou.aliyuncs.com/gitbook/qrcode.png" height="360">
    <img src="https://crawlab.oss-cn-hangzhou.aliyuncs.com/gitbook/payment.jpg" height="360">
</p>
