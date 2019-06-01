# Crawlab

![](http://114.67.75.98:8081/buildStatus/icon?job=crawlab%2Fdevelop)
![](https://img.shields.io/badge/version-v0.2.1-blue.svg)
<a href="https://github.com/tikazyq/crawlab/blob/master/LICENSE" target="_blank">
    <img src="https://img.shields.io/badge/license-BSD-blue.svg">
</a>

[中文](https://github.com/tikazyq/crawlab/blob/master/README-zh.md) | English

Celery-based web crawler admin platform for managing distributed web spiders regardless of languages and frameworks. 

[Demo](http://114.67.75.98:8080) | [Documentation](https://tikazyq.github.io/crawlab)


## Pre-requisite
- Python 3.6+
- Node.js 8.12+
- MongoDB
- Redis

## Installation

```bash
# install the requirements for backend
pip install -r requirements.txt
```

```bash
# install frontend node modules
cd frontend
npm install
```

## Configure

Please edit configuration file `config.py` to configure api and database connections.

## Quick Start
```bash
# Start backend API
python app.py

# Start Flower service
python ./bin/run_flower.py

# Start worker
python ./bin/run_worker.py
```

```bash
# run frontend client
cd frontend
npm run serve
```

## Screenshot

#### Home Page

![](https://user-gold-cdn.xitu.io/2019/3/6/169524d4c7f117f7?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

#### Spider List

![](https://user-gold-cdn.xitu.io/2019/3/6/169524daf9c8ccef?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

#### Spider Detail - Overview

![](https://user-gold-cdn.xitu.io/2019/3/6/169524e0794d6be1?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

#### Task Detail - Results

![](https://user-gold-cdn.xitu.io/2019/3/6/169524e4064c7f0a?imageView2/0/w/1280/h/960/format/webp/ignore-error/1)

## Architecture

Crawlab's architecture is very similar to Celery's, but a few more modules including Frontend, Spiders and Flower are added to feature the crawling management functionality. 

![crawlab-architecture](./docs/img/crawlab-architecture.png)

### Nodes

Nodes are actually the workers defined in Celery. A node is running and connected to a task queue, redis for example, to receive and run tasks. As spiders need to be deployed to the nodes, users should specify their ip addresses and ports before the deployment.

### Spiders

##### Auto Discovery
In `config.py` file, edit `PROJECT_SOURCE_FILE_FOLDER` as the directory where the spiders projects are located. The web app will discover spider projects automatically. How simple is that!

##### Deploy Spiders

All spiders need to be deployed to a specific node before crawling. Simply click "Deploy" button on spider detail page and the spiders will be deployed to all active nodes. 

##### Run Spiders

After deploying the spider, you can click "Run" button on spider detail page and select a specific node to start crawling. It will triggers a task for the crawling, where you can see in detail in tasks page.

### Tasks

Tasks are triggered and run by the workers. Users can view the task status, logs and results in the task detail page. 

### App

This is a Flask app that provides necessary API for common operations such as CRUD, spider deployment and task running. Each node has to run the flask app to get spiders deployed on this machine. Simply run `python manage.py app` or `python ./bin/run_app.py` to start the app.

### Broker

Broker is the same as defined in Celery. It is the queue for running async tasks.

### Frontend

Frontend is basically a Vue SPA that inherits from [Vue-Element-Admin](https://github.com/PanJiaChen/vue-element-admin) of [PanJiaChen](https://github.com/PanJiaChen). Thanks for his awesome template.

## Integration with Other Frameworks

A task is triggered via `Popen` in python `subprocess` module. A Task ID is will be defined as a variable `CRAWLAB_TASK_ID` in the shell environment to link the data to the task. 

In your spider program, you should store the `CRAWLAB_TASK_ID` value in the database with key `task_id`. Then Crawlab would know how to link those results to a particular task. For now, Crawlab only supports MongoDB. 

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

The reason is that most of the existing platforms are depending on Scrapyd, which limits the choice only within python and scrapy. Surely scrapy is a great web crawl frameowrk, but it cannot do everything. 

Crawlab is easy to use, general enough to adapt spiders in any language and any framework. It has also a beautiful frontend interface for users to manage spiders much more easily. 

|Framework | Type | Distributed | Frontend | Scrapyd-Dependent |
|:---:|:---:|:---:|:---:|:---:|
| [Crawlab](https://github.com/tikazyq/crawlab) | Admin Platform | Y | Y | N
| [Gerapy](https://github.com/Gerapy/Gerapy) | Admin Platform | Y | Y | Y
| [SpiderKeeper](https://github.com/DormyMo/SpiderKeeper) | Admin Platform | Y | Y | Y
| [ScrapydWeb](https://github.com/my8100/scrapydweb) | Admin Platform | Y | Y | Y
| [Scrapyd](https://github.com/scrapy/scrapyd) | Web Service | Y | N | N/A

## TODOs
##### Backend
- [ ] File Management
- [ ] MySQL Database Support
- [ ] Task Restart
- [ ] Node Monitoring
- [ ] More spider examples

##### Frontend
- [x] Task Stats/Analytics
- [x] Table Filters
- [x] Multi-Language Support (中文)
- [ ] Login & User Management
- [ ] General Search

## Community & Sponsorship

If you feel Crawlab could benefit your daily work or your company, please add the author's Wechat account noting "Crawlab" to enter the discussion group. Or you scan the Alipay QR code below to give us a reward to upgrade our teamwork software or buy a coffee.

<p align="center">
    <img src="https://user-gold-cdn.xitu.io/2019/3/15/169814cbd5e600e9?imageslim" height="360">
    <img src="https://raw.githubusercontent.com/tikazyq/crawlab/master/docs/img/payment.jpg" height="360">
</p>
