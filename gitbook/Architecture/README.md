## 架构

Crawlab的架构跟Celery非常相似，但是加入了包括前端、爬虫、Flower在内的额外模块，以支持爬虫管理的功能。架构图如下。

![](https://crawlab.oss-cn-hangzhou.aliyuncs.com/gitbook/architecture.png)

### 节点 Node

节点其实就是Celery中的`worker`。一个节点运行时会连接到一个任务队列（例如`Redis`）来接收和运行任务。所有爬虫需要在运行时被部署到节点上，用户在部署前需要定义节点的IP地址和端口。

### 后台应用 Backend App

这是一个Flask应用，提供了必要的API来支持常规操作，例如CRUD、爬虫部署以及任务运行。每一个节点需要启动Flask应用来支持爬虫部署。运行`python app.py`来启动应用。

### 爬虫 Spider

爬虫源代码或配置规则储存在`App`上，需要被部署到各个`worker`节点中。

### 任务 Task

任务被触发并被节点执行。用户可以在任务详情页面中看到任务到状态、日志和抓取结果。

### 中间者 Broker

中间者跟Celery中定义的一样，作为运行异步任务的队列。

### 前端 Frontend

前端其实就是一个基于[Vue-Element-Admin](https://github.com/PanJiaChen/vue-element-admin)的单页应用。其中重用了很多Element-UI的控件来支持相应的展示。

### Flower

一个Celery的插件，用于监控Celery节点。
