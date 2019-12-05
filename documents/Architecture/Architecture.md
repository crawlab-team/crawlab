## 整体架构

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
