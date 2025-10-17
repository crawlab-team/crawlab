# Crawlab

<p>
  <a href="https://github.com/crawlab-team/crawlab/actions/workflows/docker-crawlab.yml" target="_blank">
    <img src="https://github.com/crawlab-team/crawlab/workflows/Docker%20Image%20CI:%20crawlab/badge.svg">
  </a>
  <a href="https://github.com/crawlab-team/crawlab/releases" target="_blank">
    <img src="https://img.shields.io/github/release/crawlab-team/crawlab.svg?logo=github">
  </a>
  <a href="https://github.com/crawlab-team/crawlab/commits/main" target="_blank">
    <img src="https://img.shields.io/github/last-commit/crawlab-team/crawlab.svg">
  </a>
  <a href="https://github.com/crawlab-team/crawlab/issues?q=is%3Aissue+is%3Aopen+label%3Abug" target="_blank">
    <img src="https://img.shields.io/github/issues/crawlab-team/crawlab/bug.svg?label=bugs&color=red">
  </a>
  <a href="https://github.com/crawlab-team/crawlab/issues?q=is%3Aissue+is%3Aopen+label%3Aenhancement" target="_blank">
    <img src="https://img.shields.io/github/issues/crawlab-team/crawlab/enhancement.svg?label=enhancements&color=cyan">
  </a>
  <a href="https://github.com/crawlab-team/crawlab/blob/main/LICENSE" target="_blank">
    <img src="https://img.shields.io/github/license/crawlab-team/crawlab.svg">
  </a>
</p>

中文 | [English](https://github.com/crawlab-team/crawlab/blob/main/README.md)

[安装](#安装) | [运行](#运行) | [截图](#截图) | [架构](#架构) | [集成](#与其他框架的集成) | [比较](#与其他框架比较) | [社区与赞助](#社区与赞助) | [更新日志](https://github.com/crawlab-team/crawlab/blob/main/CHANGELOG.md) | [免责声明](https://github.com/crawlab-team/crawlab/blob/main/DISCLAIMER.md)

基于 Golang 的分布式网络爬虫管理平台，支持多种编程语言，包括 Python、NodeJS、Go、Java，以及多种网络爬虫框架，包括 Scrapy、Puppeteer、Selenium。

[在线演示](https://demo.crawlab.cn) | [文档](https://docs.crawlab.cn)


## 安装

您可以参考[安装指南](https://docs.crawlab.cn/getting-started/installation)。

## 快速开始

请打开命令行提示符并执行以下命令。请确保您已提前安装了 [Docker](https://www.docker.com)。

```bash
git clone https://github.com/crawlab-team/examples
cd examples/docker/basic
docker-compose up -d
```

接下来，您可以查看 `docker-compose.yml`（包含详细的配置参数）和[文档](http://docs.crawlab.cn)以获取更多信息。

## 运行

### Docker

请使用 `docker compose` 一键启动。这样，您甚至不需要配置 MongoDB 数据库。创建一个名为 `docker-compose.yml` 的文件并输入以下内容。

```yaml
services:
  master: 
    image: crawlabteam/crawlab:latest
    container_name: crawlab_example_master
    environment:
      CRAWLAB_NODE_MASTER: "Y"
      CRAWLAB_MONGO_HOST: "mongo"
    volumes:
      - "./.crawlab/master:/root/.crawlab"
    ports:    
      - "8080:8080"
    depends_on:
      - mongo

  worker01: 
    image: crawlabteam/crawlab:latest
    container_name: crawlab_example_worker01
    environment:
      CRAWLAB_NODE_MASTER: "N"
      CRAWLAB_MASTER_HOST: "master"
    volumes:
      - "./.crawlab/worker01:/root/.crawlab"
    depends_on:
      - master

  worker02: 
    image: crawlabteam/crawlab:latest
    container_name: crawlab_example_worker02
    environment:
      CRAWLAB_NODE_MASTER: "N"
      CRAWLAB_GRPC_ADDRESS: "master"
      CRAWLAB_FS_FILER_URL: "http://master:8080/api/filer"
    volumes:
      - "./.crawlab/worker02:/root/.crawlab"
    depends_on:
      - master

  mongo:
    image: mongo:4.2
    container_name: crawlab_example_mongo
    restart: always
```

然后执行以下命令，Crawlab 主节点和工作节点 + MongoDB 将启动。打开浏览器并输入 `http://localhost:8080` 即可看到用户界面。

```bash
docker-compose up -d
```

有关 Docker 部署的详细信息，请参阅[相关文档](https://docs.crawlab.cn/en/guide/installation/docker.html)。



## 截图

#### 登录

![]( https://github.com/crawlab-team/images/blob/main/20210729/screenshot-login.png?raw=true)

#### 主页

![]( https://github.com/crawlab-team/images/blob/main/20210729/screenshot-home.png?raw=true)

#### 节点列表

![]( https://github.com/crawlab-team/images/blob/main/20210729/screenshot-node-list.png?raw=true)

#### 爬虫列表

![](https://github.com/crawlab-team/images/blob/main/20210729/screenshot-spider-list.png?raw=true)

#### 爬虫概览

![](https://github.com/crawlab-team/images/blob/main/20210729/screenshot-spider-detail-overview.png?raw=true)

#### 爬虫文件

![](https://github.com/crawlab-team/images/blob/main/20210729/screenshot-spider-detail-files.png?raw=true)

#### 任务日志

![](https://github.com/crawlab-team/images/blob/main/20210729/screenshot-task-detail-logs.png?raw=true)

#### 任务结果

![](https://github.com/crawlab-team/images/blob/main/20210729/screenshot-task-detail-data.png?raw=true)

#### 定时任务

![](https://github.com/crawlab-team/images/blob/main/20210729/screenshot-schedule-detail-overview.png?raw=true)



## 架构

Crawlab 的架构由一个主节点、多个工作节点、[SeaweedFS](https://github.com/chrislusf/seaweedfs)（分布式文件系统）和 MongoDB 数据库组成。

![](https://github.com/crawlab-team/images/blob/main/20210729/crawlab-architecture-v0.6.png?raw=true)

前端应用与主节点交互，主节点与其他组件通信，如 MongoDB、SeaweedFS 和工作节点。主节点和工作节点通过 [gRPC](https://grpc.io)（一个 RPC 框架）相互通信。任务由主节点中的任务调度器模块调度，并由工作节点中的任务处理器模块接收，在任务执行器中执行这些任务。任务执行器实际上是运行爬虫或爬取程序的进程，也可以通过 gRPC（集成在 SDK 中）将数据发送到其他数据源，例如 MongoDB。

### 主节点

主节点是 Crawlab 架构的核心。它是 Crawlab 的中央控制系统。

主节点提供以下服务：
1. 任务调度；
2. 工作节点管理和通信；
3. 爬虫部署；
4. 前端和 API 服务；
5. 任务执行（您可以将主节点视为工作节点）

主节点与前端应用通信，并将爬取任务发送给工作节点。同时，主节点将爬虫上传（部署）到分布式文件系统 SeaweedFS，供工作节点同步。

### 工作节点

工作节点的主要功能是执行爬取任务并存储结果和日志，通过 gRPC 与主节点通信。通过增加工作节点的数量，Crawlab 可以横向扩展，不同的爬取任务可以分配给不同的节点执行。

### MongoDB

MongoDB 是 Crawlab 的运营数据库。它存储节点、爬虫、任务、调度等数据。任务队列也存储在 MongoDB 中。

### SeaweedFS

SeaweedFS 是由 [Chris Lu](https://github.com/chrislusf) 创建的开源分布式文件系统。它可以在分布式系统中稳健地存储和共享文件。在 Crawlab 中，SeaweedFS 主要作为文件同步系统和存储任务日志文件的地方。

### 前端

前端应用基于 [Element-Plus](https://github.com/element-plus/element-plus) 构建，这是一个流行的基于 [Vue 3](https://github.com/vuejs/vue-next) 的 UI 框架。它与主节点上托管的 API 交互，并间接控制工作节点。



## 与其他框架的集成

[Crawlab SDK](https://github.com/crawlab-team/crawlab-sdk) 提供了一些 `helper` 方法，使您更容易将爬虫集成到 Crawlab 中，例如保存结果。

### Scrapy

在 Scrapy 项目的 `settings.py` 中，找到名为 `ITEM_PIPELINES` 的变量（一个 `dict` 类型变量）。添加以下内容。

```python
ITEM_PIPELINES = {
    'crawlab.scrapy.pipelines.CrawlabPipeline': 888,
}
```

然后，启动 Scrapy 爬虫。完成后，您应该能够在**任务详情 -> 数据**中看到抓取的结果。

### 通用 Python 爬虫

请将以下内容添加到您的爬虫文件中以保存结果。

```python
# 导入结果保存方法
from crawlab import save_item

# 这是一个结果记录，必须是 dict 类型
result = {'name': 'crawlab'}

# 调用结果保存方法
save_item(result)
```

然后，启动爬虫。完成后，您应该能够在**任务详情 -> 数据**中看到抓取的结果。

### 其他框架 / 语言

爬取任务实际上是通过 shell 命令执行的。任务 ID 将以名为 `CRAWLAB_TASK_ID` 的环境变量的形式传递给爬取任务进程。通过这样做，数据可以与任务关联。



## 与其他框架比较

目前已经存在一些爬虫管理框架了。那么为什么要使用 Crawlab？

原因是大多数现有平台都依赖于 Scrapyd，这将选择限制在 Python 和 Scrapy 范围内。当然，Scrapy 是一个很棒的网络爬取框架，但它不能做所有事情。

Crawlab 易于使用，足够通用，可以适应任何语言和任何框架的爬虫。它还有一个漂亮的前端界面，让用户可以更轻松地管理爬虫。

|框架 | 技术 | 优点 | 缺点 | Github 统计数据 |
|:---|:---|:---|-----| :---- |
| [Crawlab](https://github.com/crawlab-team/crawlab) | Golang + Vue|不局限于 Scrapy，适用于所有编程语言和框架。精美的 UI 界面。天然支持分布式爬虫。支持爬虫管理、任务管理、定时任务、结果导出、分析、通知、可配置爬虫、在线代码编辑器等。|暂不支持爬虫版本控制| ![](https://img.shields.io/github/stars/crawlab-team/crawlab) ![](https://img.shields.io/github/forks/crawlab-team/crawlab) |
| [ScrapydWeb](https://github.com/my8100/scrapydweb) | Python Flask + Vue|精美的 UI 界面，内置 Scrapy 日志解析器，任务执行的统计和图表，支持节点管理、定时任务、邮件通知、移动端。功能齐全的爬虫管理平台。|不支持 Scrapy 以外的爬虫。由于 Python Flask 后端的限制，性能有限。| ![](https://img.shields.io/github/stars/my8100/scrapydweb) ![](https://img.shields.io/github/forks/my8100/scrapydweb) |
| [Gerapy](https://github.com/Gerapy/Gerapy) | Python Django + Vue|Gerapy 由网络爬虫大师 [Germey Cui](https://github.com/Germey) 构建。安装和部署简单。精美的 UI 界面。支持节点管理、代码编辑、可配置的爬取规则等。|同样不支持 Scrapy 以外的爬虫。根据用户反馈，v1.0 存在很多 bug。期待 v2.0 的改进| ![](https://img.shields.io/github/stars/Gerapy/Gerapy) ![](https://img.shields.io/github/forks/Gerapy/Gerapy) |
| [SpiderKeeper](https://github.com/DormyMo/SpiderKeeper) | Python Flask|开源的 Scrapyhub。简洁的 UI 界面。支持定时任务。|可能过于简化，不支持分页，不支持节点管理，不支持 Scrapy 以外的爬虫。| ![](https://img.shields.io/github/stars/DormyMo/SpiderKeeper) ![](https://img.shields.io/github/forks/DormyMo/SpiderKeeper) |



## 贡献者
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
<a href="https://github.com/zkqiang">
  <img src="https://avatars3.githubusercontent.com/u/32983588?s=460&u=83082ddc0a3020279374b94cce70f1aebb220b3d&v=4" height="80">
</a>

## JetBrains 支持

<p align="center">
  <a href="https://www.jetbrains.com" target="_blank">
    <img src="https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.png" height="360">
  </a>
</p>

## 社区

如果您觉得 Crawlab 对您的日常工作或公司有帮助，请添加作者的微信账号并注明"Crawlab"以进入讨论群。

<p align="center">
    <img src="https://crawlab.oss-cn-hangzhou.aliyuncs.com/gitbook/qrcode.png" height="360">
</p>
