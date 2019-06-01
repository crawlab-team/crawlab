# Crawlab

![](http://114.67.75.98:8081/buildStatus/icon?job=crawlab%2Fdevelop)
![](https://img.shields.io/badge/版本-v0.2.3-blue.svg)
<a href="https://github.com/tikazyq/crawlab/blob/master/LICENSE" target="_blank">
    <img src="https://img.shields.io/badge/License-BSD-blue.svg">
</a>

中文 | [English](https://github.com/tikazyq/crawlab/blob/master/README.md)

基于Celery的爬虫分布式爬虫管理平台，支持多种编程语言以及多种爬虫框架.

[查看演示 Demo](http://114.67.75.98:8080) | [文档](https://tikazyq.github.io/crawlab)

## 要求
- Python 3.6+
- Node.js 8.12+
- MongoDB
- Redis

## 安装

```bash
# 安装后台类库
pip install -r requirements.txt
```

```bash
# 安装前台类库
cd frontend
npm install
```

## 配置

请更改配置文件`config.py`，配置API和数据库连接.

## 快速开始
```bash
python manage.py serve
```

## 截图

#### 首页

![](https://user-gold-cdn.xitu.io/2019/3/6/169524d4c7f117f7?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

#### 爬虫列表

![](https://user-gold-cdn.xitu.io/2019/3/6/169524daf9c8ccef?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

#### 爬虫详情 - 概览

![](https://user-gold-cdn.xitu.io/2019/3/6/169524e0794d6be1?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

#### 任务详情 - 抓取结果

![](https://user-gold-cdn.xitu.io/2019/3/6/169524e4064c7f0a?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

## 使用流程

![user-process](./docs/img/用户使用流程图.png)

## 架构

Crawlab的架构跟Celery非常相似，但是加入了包括前端、爬虫、Flower在内的额外模块，以支持爬虫管理的功能。

![crawlab-architecture](./docs/img/crawlab-architecture.png)

### 节点

节点其实就是Celery中的Worker。一个节点运行时会连接到一个任务队列（例如Redis）来接收和运行任务。所有爬虫需要在运行时被部署到节点上，用户在部署前需要定义节点的IP地址和端口。

### 爬虫

##### 自动发现

在`config.py`文件中，修改变量`PROJECT_SOURCE_FILE_FOLDER`作为爬虫项目所在的目录。Crawlab后台程序会自动发现这些爬虫项目并储存到数据库中。是不是很方便？

##### 部署爬虫

所有爬虫需要在抓取前被部署当相应当节点中。在"爬虫详情"页面点击"Deploy"按钮，爬虫将被部署到所有有效到节点中。

##### 运行爬虫

部署爬虫之后，你可以在"爬虫详情"页面点击"Run"按钮来启动爬虫。一个爬虫任务将被触发，你可以在任务列表页面中看到这个任务。

### 任务

任务被触发并被节点执行。用户可以在任务详情页面中看到任务到状态、日志和抓取结果。

### 后台应用

这是一个Flask应用，提供了必要的API来支持常规操作，例如CRUD、爬虫部署以及任务运行。每一个节点需要启动Flask应用来支持爬虫部署。运行`python manage.py app`或`python ./bin/run_app.py`来启动应用。

### 中间者

中间者跟Celery中定义的一样，作为运行异步任务的队列。

### 前端

前端其实就是一个基于[Vue-Element-Admin](https://github.com/PanJiaChen/vue-element-admin)的单页应用。其中重用了很多Element-UI的控件来支持相应的展示。

## 与其他框架的集成

任务是利用python的`subprocess`模块中的`Popen`来实现的。任务ID将以环境变量`CRAWLAB_TASK_ID`的形式存在于爬虫任务运行的进程中，并以此来关联抓取数据。

在你的爬虫程序中，你需要将`CRAWLAB_TASK_ID`的值以`task_id`作为可以存入数据库中。这样Crawlab就直到如何将爬虫任务与抓取数据关联起来了。当前，Crawlab只支持MongoDB。

### Scrapy

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
| [Gerapy](https://github.com/Gerapy/Gerapy) | 管理平台 | Y | Y | Y
| [SpiderKeeper](https://github.com/DormyMo/SpiderKeeper) | 管理平台 | Y | Y | Y
| [ScrapydWeb](https://github.com/my8100/scrapydweb) | 管理平台 | Y | Y | Y
| [Scrapyd](https://github.com/scrapy/scrapyd) | 网络服务 | Y | N | N/A

## TODOs
##### 后端
- [ ] 文件管理
- [ ] MySQL数据库支持
- [ ] 重跑任务
- [ ] 节点监控
- [ ] 更多爬虫例子

##### 前端
- [x] 任务数据统计
- [x] 表格过滤
- [x] 多语言支持 (中文)
- [ ] 登录和用户管理
- [ ] 全局搜索

## 社区 & 赞助

如果您觉得Crawlab对您的日常开发或公司有帮助，请加作者微信 tikazyq1 并注明"Crawlab"，作者会将你拉入群。或者，您可以扫下方支付宝二维码给作者打赏去升级团队协作软件或买一杯咖啡。

<p align="center">
    <img src="https://user-gold-cdn.xitu.io/2019/3/15/169814cbd5e600e9?imageslim" height="360">
    <img src="https://raw.githubusercontent.com/tikazyq/crawlab/master/docs/img/payment.jpg" height="360">
</p>
