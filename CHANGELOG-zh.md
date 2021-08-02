# 0.6.0 (TBC)

(TBC)

# 0.5.1 (2020-07-31)

### 功能 / 优化
- **加入错误详情信息**.
- **加入 Golang 编程语言支持**.
- **加入 Chrome Driver 和 Firefox 的 Web Driver 安装脚本**.
- **支持系统任务**. "系统任务"跟普通爬虫任务相似，允许用户查看诸如安装语言之类的任务日志.
- **将安装语言从 RPC 更改为系统任务**.

### Bug 修复
- **修复在爬虫市场中第一次下载爬虫时会报500错误**. [#808](https://github.com/crawlab-team/crawlab/issues/808)
- **修复一部分翻译问题**.
- **修复任务详情 500 错误**. [#810](https://github.com/crawlab-team/crawlab/issues/810)
- **修复密码重置问题**. [#811](https://github.com/crawlab-team/crawlab/issues/811)
- **修复无法下载 CSV 问题**. [#812](https://github.com/crawlab-team/crawlab/issues/812)
- **修复无法安装 Node.js 问题**. [#813](https://github.com/crawlab-team/crawlab/issues/813)
- **修复批量添加定时任务时默认为禁用问题**. [#814](https://github.com/crawlab-team/crawlab/issues/814)

# 0.5.0 (2020-07-19)
### 功能 / 优化
- **爬虫市场**. 允许用户下载开源爬虫到 Crawlab.
- **批量操作**. 允许用户与 Crawlab 批量交互，例如批量运行任务、批量删除爬虫等等.
- **迁移 MongoDB 驱动器至 `MongoDriver`**.
- **重构优化节点逻辑代码**.
- **更改默认 `task.workers` 至 16**.
- **更改默认 nginx `client_max_body_size` 为 200m**.
- **支持写日志到 ElasticSearch**.
- **在 Scrapy 页面展示错误详情**.
- **删除挑战页面**.
- **将反馈、免责声明页面移动到顶部**.

### Bug 修复
- **修复由于 TTL 索引未创建导致的日志不过期问题**.
- **设置默认日志过期时间为 1 天**.
- **`task_id` 索引没有创建**.
- **`docker-compose.yml` 修复**.
- **修复 404 页面**.
- **修复无法先创建工作节点问题**.

# 0.4.10 (2020-04-21)
### 功能 / 优化
- **优化日志管理**. 集中化管理日志，储存在 MongoDB，减少对 PubSub 的依赖，允许日志异常检测.
- **自动安装依赖**. 允许从 `requirements.txt` 和 `package.json` 自动安装依赖.
- **API Token**. 允许用户生成 API Token，并利用它们来集成到自己的系统中.
- **Web Hook**. 当任务开始或结束时，触发 Web Hook http 请求到预定义好的 URL.
- **自动生成结果集**. 如果没有设置，自动设置结果集为 `results_<spider_name>`.
- **优化项目列表**. 项目列表中不展示 "No Project".
- **升级 Node.js**. 将 Node.js 版本从 v8.12 升级到 v10.19.
- **定时任务增加运行按钮**. 允许用户在定时任务界面手动运行爬虫任务.

### Bug 修复
- **无法注册**. [#670](https://github.com/crawlab-team/crawlab/issues/670)
- **爬虫定时任务标签 Cron 表达式显示秒**. [#678](https://github.com/crawlab-team/crawlab/issues/678)
- **爬虫每日数据缺失**. [#684](https://github.com/crawlab-team/crawlab/issues/684)
- **结果数量未即时更新**. [#689](https://github.com/crawlab-team/crawlab/issues/689)

# 0.4.9 (2020-03-31)
### 功能 / 优化
- **挑战**. 用户可以完成不同的趣味挑战..
- **更高级的权限控制**. 更细化的权限管理，例如普通用户只能查看或管理自己的爬虫或项目，而管理用户可以查看或管理所有爬虫或项目.
- **反馈**. 允许用户发送反馈和评分给 Crawlab 开发组.
- **更好的主页指标**. 优化主页上的指标展示.
- **可配置爬虫转化为自定义爬虫**. 用户可以将自己的可配置爬虫转化为 Scrapy 自定义爬虫.
- **查看定时任务触发的任务**. 允许用户查看定时任务触发的任务. [#648](https://github.com/crawlab-team/crawlab/issues/648)
- **支持结果去重**. 允许用户配置结果去重. [#579](https://github.com/crawlab-team/crawlab/issues/579)
- **支持任务重试**. 允许任务重新触发历史任务.

### Bug 修复
- **无法注册**. [#670](https://github.com/crawlab-team/crawlab/issues/670)
- **CLI 无法在 Windows 上使用**. [#580](https://github.com/crawlab-team/crawlab/issues/580)
- **重新上传错误**. [#643](https://github.com/crawlab-team/crawlab/issues/643) [#640](https://github.com/crawlab-team/crawlab/issues/640)
- **上传丢失文件目录**. [#646](https://github.com/crawlab-team/crawlab/issues/646)
- **无法在爬虫定时任务标签中添加定时任务**.

# 0.4.8 (2020-03-11)
### 功能 / 优化
- **支持更多编程语言安装**. 现在用户可以安装或预装更多的编程语言，包括 Java、.Net Core、PHP.
- **安装 UI 优化**. 用户能够更好的查看和管理节点列表页的安装.
- **更多 Git 支持**. 允许用户查看 Git Commits 记录，并 Checkout 到相应 Commit.
- **支持用 Hostname 作为节点注册类型**. 用户可以将 hostname 作为节点的唯一识别号.
- **RPC 支持**. 加入 RPC 支持来更好的管理节点通信.
- **是否在主节点运行开关**. 用户可以决定是否在主节点运行，如果为否，则所有任务将在工作节点上运行.
- **默认禁用教程**.
- **加入相关文档侧边栏**.
- **加载页面优化**.

### Bug 修复
- **重复节点**. [#391](https://github.com/crawlab-team/crawlab/issues/391)
- **重复上传爬虫**. [#603](https://github.com/crawlab-team/crawlab/issues/603)
- **节点第三方模块安装失败导致 节点安装第三方部分无法使用**. [#609](https://github.com/crawlab-team/crawlab/issues/609)
- **离线节点也会创建任务**. [#622](https://github.com/crawlab-team/crawlab/issues/622)

# 0.4.7 (2020-02-24)
### 功能 / 优化
- **更好的支持 Scrapy**. 爬虫识别，`settings.py` 配置，日志级别选择，爬虫选择. [#435](https://github.com/crawlab-team/crawlab/issues/435)
- **Git 同步**. 允许用户将 Git 项目同步到 Crawlab.
- **长任务支持**. 用户可以添加长任务爬虫，这些爬虫可以跑长期运行的任务. [425](https://github.com/crawlab-team/crawlab/issues/425)
- **爬虫列表优化**. 分状态任务列数统计，任务列表详情弹出框，图例. [425](https://github.com/crawlab-team/crawlab/issues/425)
- **版本升级检测**. 检测最新版本，通知用户升级.
- **批量操作爬虫**. 允许用户批量运行/停止爬虫任务，以及批量删除爬虫.
- **复制爬虫**. 允许用户复制已存在爬虫来创建新爬虫.
- **微信群二维码**.

### Bug 修复
- **定时任务爬虫选择问题**. 字段不会随着爬虫变化而响应.
- **定时任务冲突问题**. 两个不同的爬虫设置定时任务，时间设置成相同的话，可能会有bug. [#515](https://github.com/crawlab-team/crawlab/issues/515) [#565](https://github.com/crawlab-team/crawlab/issues/565)
- **任务日志问题**. 在同一时间触发的不同任务可能会写入同一个日志文件. [#577](https://github.com/crawlab-team/crawlab/issues/577)
- **任务列表筛选选项不全**.

# 0.4.6 (2020-02-13)
### 功能 / 优化
- **Node.js SDK**. 用户可以将 SDK 应用到他们的 Node.js 爬虫中.
- **日志管理优化**. 日志搜索，错误高亮，自动滚动.
- **任务执行流程优化**. 允许用户在触发任务后跳转到该任务详情页.
- **任务展示优化**. 在爬虫详情页的最近任务表格中加入了“参数”列. [#295](https://github.com/crawlab-team/crawlab/issues/295)
- **爬虫列表优化**. 在爬虫列表页加入"更新时间"和"创建时间". [#505](https://github.com/crawlab-team/crawlab/issues/505)
- **页面加载占位器**. 

### Bug 修复
- **定时任务配置失去焦点**. [#519](https://github.com/crawlab-team/crawlab/issues/519)
- **无法用 CLI 工具上传爬虫**. [#524](https://github.com/crawlab-team/crawlab/issues/524)

# 0.4.5 (2020-02-03)
### 功能 / 优化
- **交互式教程**. 引导用户了解 Crawlab 的主要功能.
- **加入全局环境变量**. 可以设置全局环境变量，然后传入到所有爬虫程序中. [#177](https://github.com/crawlab-team/crawlab/issues/177)
- **项目**. 允许用户将爬虫关联到项目上. [#316](https://github.com/crawlab-team/crawlab/issues/316)
- **示例爬虫**. 当初始化时，自动加入示例爬虫. [#379](https://github.com/crawlab-team/crawlab/issues/379)
- **用户管理优化**. 限制管理用户的权限. [#456](https://github.com/crawlab-team/crawlab/issues/456)
- **设置页面优化**.
- **任务结果页面优化**.

### Bug 修复
- **无法找到爬虫文件错误**. [#485](https://github.com/crawlab-team/crawlab/issues/485)
- **点击删除按钮导致跳转**. [#480](https://github.com/crawlab-team/crawlab/issues/480)
- **无法在空爬虫里创建文件**. [#479](https://github.com/crawlab-team/crawlab/issues/479)
- **下载结果错误**. [#465](https://github.com/crawlab-team/crawlab/issues/465)
- **crawlab-sdk CLI 错误**. [#458](https://github.com/crawlab-team/crawlab/issues/458)
- **页面刷新问题**. [#441](https://github.com/crawlab-team/crawlab/issues/441)
- **结果不支持 JSON**. [#202](https://github.com/crawlab-team/crawlab/issues/202)
- **修复“删除爬虫后获取所有爬虫”错误**.
- **修复 i18n 警告**.

# 0.4.4 (2020-01-17)

### 功能 / 优化
- **邮件通知**. 允许用户发送邮件消息通知.
- **钉钉机器人通知**. 允许用户发送钉钉机器人消息通知.
- **企业微信机器人通知**. 允许用户发送企业微信机器人消息通知.
- **API 地址优化**. 在前端加入相对路径，因此用户不需要特别注明 `CRAWLAB_API_ADDRESS`.
- **SDK 兼容**. 允许用户通过 Crawlab SDK 与 Scrapy 或通用爬虫集成.
- **优化文件管理**. 加入树状文件侧边栏，让用户更方便的编辑文件.
- **高级定时任务 Cron**. 允许用户通过 Cron 可视化编辑器编辑定时任务.

### Bug 修复
- **`nil retuened` 错误**.
- **使用 HTTPS 出现的报错**.
- **无法在爬虫列表页运行可配置爬虫**.
- **上传爬虫文件缺少表单验证**.

# 0.4.3 (2020-01-07)

### 功能 / 优化
- **依赖安装**. 允许用户在平台 Web 界面安装/卸载依赖以及添加编程语言（暂时只有 Node.js）。
- **Docker 中预装编程语言**. 允许 Docker 用户通过设置 `CRAWLAB_SERVER_LANG_NODE` 为 `Y` 来预装 `Node.js` 环境.
- **在爬虫详情页添加定时任务列表**. 允许用户在爬虫详情页查看、添加、编辑定时任务. [#360](https://github.com/crawlab-team/crawlab/issues/360)
- **Cron 表达式与 Linux 一致**. 将表达式从 6 元素改为 5 元素，与 Linux 一致.
- **启用/禁用定时任务**. 允许用户启用/禁用定时任务. [#297](https://github.com/crawlab-team/crawlab/issues/297)
- **优化任务管理**. 允许用户批量删除任务. [#341](https://github.com/crawlab-team/crawlab/issues/341)
- **优化爬虫管理**. 允许用户在爬虫列表页对爬虫进行筛选和排序.
- **添加中文版 `CHANGELOG`**.
- **在顶部添加 Github 加星按钮**.

### Bug 修复
- **定时任务问题**. [#423](https://github.com/crawlab-team/crawlab/issues/423)
- **上传爬虫zip文件问题**. [#403](https://github.com/crawlab-team/crawlab/issues/403) [#407](https://github.com/crawlab-team/crawlab/issues/407)
- **因为网络原因导致崩溃**. [#340](https://github.com/crawlab-team/crawlab/issues/340)
- **定时任务无法正常运行**
- **定时任务列表列表错位问题**
- **刷新按钮跳转错误问题**

# 0.4.2 (2019-12-26)
### 功能 / 优化
- **免责声明**. 加入免责声明.
- **通过 API 获取版本号**. [#371](https://github.com/crawlab-team/crawlab/issues/371)
- **通过配置来允许用户注册**. [#346](https://github.com/crawlab-team/crawlab/issues/346)
- **允许添加新用户**.
- **更高级的文件管理**. 允许用户添加、编辑、重命名、删除代码文件. [#286](https://github.com/crawlab-team/crawlab/issues/286)
- **优化爬虫创建流程**. 允许用户在上传 zip 文件前创建空的自定义爬虫.
- **优化任务管理**. 允许用户通过选择条件过滤任务. [#341](https://github.com/crawlab-team/crawlab/issues/341)

### Bug 修复
- **重复节点**. [#391](https://github.com/crawlab-team/crawlab/issues/391)
- **"mongodb no reachable" 错误**. [#373](https://github.com/crawlab-team/crawlab/issues/373)

# 0.4.1 (2019-12-13)
### 功能 / 优化
- **Spiderfile 优化**. 将阶段由数组更换为字典. [#358](https://github.com/crawlab-team/crawlab/issues/358)
- **百度统计更新**.

### Bug 修复
- **无法展示定时任务**. [#353](https://github.com/crawlab-team/crawlab/issues/353)
- **重复节点注册**. [#334](https://github.com/crawlab-team/crawlab/issues/334)

# 0.4.0 (2019-12-06)
### 功能 / 优化
- **可配置爬虫**. 允许用户添加 `Spiderfile` 来配置抓取规则.
- **执行模式**. 允许用户选择 3 种任务执行模式: *所有节点*, *指定节点* and *随机*.

### Bug 修复
- **任务意外被杀死**. [#306](https://github.com/crawlab-team/crawlab/issues/306)
- **文档更正**. [#301](https://github.com/crawlab-team/crawlab/issues/258) [#301](https://github.com/crawlab-team/crawlab/issues/258)
- **直接部署与 Windows 不兼容**. [#288](https://github.com/crawlab-team/crawlab/issues/288)
- **日志文件丢失**. [#269](https://github.com/crawlab-team/crawlab/issues/269)

# 0.3.5 (2019-10-28)
### 功能 / 优化
- **优雅关闭**. [详情](https://github.com/crawlab-team/crawlab/commit/63fab3917b5a29fd9770f9f51f1572b9f0420385)
- **节点信息优化**. [详情](https://github.com/crawlab-team/crawlab/commit/973251a0fbe7a2184ac0da09e0404a17c736aee7)
- **将系统环境变量添加到任务**. [详情](https://github.com/crawlab-team/crawlab/commit/4ab4892471965d6342d30385578ca60dc51f8ad3)
- **自动刷新任务日志**. [详情](https://github.com/crawlab-team/crawlab/commit/4ab4892471965d6342d30385578ca60dc51f8ad3)
- **允许 HTTPS 部署**. [详情](https://github.com/crawlab-team/crawlab/commit/5d8f6f0c56768a6e58f5e46cbf5adff8c7819228)

### Bug 修复
- **定时任务中无法获取爬虫列表**. [详情](https://github.com/crawlab-team/crawlab/commit/311f72da19094e3fa05ab4af49812f58843d8d93)
- **无法获取工作节点信息**. [详情](https://github.com/crawlab-team/crawlab/commit/6af06efc17685a9e232e8c2b5fd819ec7d2d1674)
- **运行爬虫任务时无法选择节点**. [详情](https://github.com/crawlab-team/crawlab/commit/31f8e03234426e97aed9b0bce6a50562f957edad)
- **结果量很大时无法获取结果数量**. [#260](https://github.com/crawlab-team/crawlab/issues/260)
- **定时任务中的节点问题**. [#244](https://github.com/crawlab-team/crawlab/issues/244)


# 0.3.1 (2019-08-25)
### 功能 / 优化
- **Docker 镜像优化**. 将 Docker 镜像进一步分割成 alpine 镜像版本的 master、worker、frontendSplit docker further into master, worker, frontend.
- **单元测试**. 用单元测试覆盖部分后端代码.
- **前端优化**. 登录页、按钮大小、上传 UI 提示. 
- **更灵活的节点注册**. 允许用户传一个变量作为注册 key，而不是默认的 MAC 地址.

### Bug 修复
- **上传大爬虫文件错误**. 上传大爬虫文件时的内存崩溃问题. [#150](https://github.com/crawlab-team/crawlab/issues/150)
- **无法同步爬虫**. 通过提高写权限等级来修复同步爬虫文件时的问题. [#114](https://github.com/crawlab-team/crawlab/issues/114)
- **爬虫页问题**. 通过删除 `Site` 字段来修复. [#112](https://github.com/crawlab-team/crawlab/issues/112)
- **节点展示问题**. 当在多个机器上跑 Docker 容器时，节点无法正确展示. [#99](https://github.com/crawlab-team/crawlab/issues/99)

# 0.3.0 (2019-07-31)
### 功能 / 优化
- **Golang 后端**: 将后端由 Python 重构为 Golang，很大的提高了稳定性和性能.
- **节点网络图**: 节点拓扑图可视化.
- **节点系统信息**: 可以查看包括操作系统、CPU数量、可执行文件在内的系统信息.
- **节点监控改进**: 节点通过 Redis 来监控和注册.
- **文件管理**: 可以在线编辑爬虫文件，包括代码高亮.
- **登录页/注册页/用户管理**: 要求用户登录后才能使用 Crawlab, 允许用户注册和用户管理，有一些基于角色的鉴权机制.
- **自动部署爬虫**: 爬虫将被自动部署或同步到所有在线节点.
- **更小的 Docker 镜像**: 瘦身版 Docker 镜像，通过多阶段构建将 Docker 镜像大小从 1.3G 减小到 700M 左右.

### Bug 修复
- **节点状态**. 节点状态不会随着节点下线而更新. [#87](https://github.com/tikazyq/crawlab/issues/87) 
- **爬虫部署错误**. 通过自动爬虫部署来修复 [#83](https://github.com/tikazyq/crawlab/issues/83) 
- **节点无法显示**. 节点无法显示在线 [#81](https://github.com/tikazyq/crawlab/issues/81) 
- **定时任务无法工作**. 通过 Golang 后端修复 [#64](https://github.com/tikazyq/crawlab/issues/64) 
- **Flower 错误**. 通过 Golang 后端修复 [#57](https://github.com/tikazyq/crawlab/issues/57) 

# 0.2.4 (2019-07-07)
### 功能 / 优化
- **文档**: 更优和更详细的文档.
- **更好的 Crontab**: 通过 UI 界面生成 Cron 表达式.
- **更优的性能**: 从原生 flask 引擎 切换到 `gunicorn`. [#78](https://github.com/tikazyq/crawlab/issues/78)

### Bug 修复
- **删除爬虫**. 删除爬虫时不止在数据库中删除，还应该删除相关的文件夹、任务和定时任务. [#69](https://github.com/tikazyq/crawlab/issues/69)
- **MongoDB 授权**. 允许用户注明 `authenticationDatabase` 来连接 `mongodb`. [#68](https://github.com/tikazyq/crawlab/issues/68)
- **Windows 兼容性**. 加入 `eventlet` 到 `requirements.txt`. [#59](https://github.com/tikazyq/crawlab/issues/59)


# 0.2.3 (2019-06-12)
### 功能 / 优化
- **Docker**: 用户能够运行 Docker 镜像来加快部署.
- **CLI**: 允许用户通过命令行来执行 Crawlab 程序.
- **上传爬虫**: 允许用户上传自定义爬虫到 Crawlab.
- **预览时编辑字段**: 允许用户在可配置爬虫中预览数据时编辑字段.

### Bug 修复
- **爬虫分页**. 爬虫列表页中修复分页问题.

# 0.2.2 (2019-05-30)
### 功能 / 优化
- **自动抓取字段**: 在可配置爬虫列表页种自动抓取字段.
- **下载结果**: 允许下载结果为 CSV 文件.
- **百度统计**: 允许用户选择是否允许向百度统计发送统计数据.

### Bug 修复
- **结果页分页**. [#45](https://github.com/tikazyq/crawlab/issues/45)
- **定时任务重复触发**: 将 Flask DEBUG 设置为 False 来保证定时任务无法重复触发. [#32](https://github.com/tikazyq/crawlab/issues/32)
- **前端环境**: 添加 `VUE_APP_BASE_URL` 作为生产环境模式变量，然后 API 不会永远都是 `localhost` [#30](https://github.com/tikazyq/crawlab/issues/30)

# 0.2.1 (2019-05-27)
- **可配置爬虫**: 允许用户创建爬虫来抓取数据，而不用编写代码.

# 0.2 (2019-05-10)

- **高级数据统计**: 爬虫详情页的高级数据统计.
- **网站数据**: 加入网站列表（中国），允许用户查看 robots.txt、首页响应时间等信息.

# 0.1.1 (2019-04-23)

- **基础统计**: 用户可以查看基础统计数据，包括爬虫和任务页中的失败任务数、结果数.
- **近实时任务信息**: 周期性（5 秒）向服务器轮训数据来实现近实时查看任务信息.
- **定时任务**: 利用 apscheduler 实现定时任务，允许用户设置类似 Cron 的定时任务.

# 0.1 (2019-04-17)

- **首次发布**
