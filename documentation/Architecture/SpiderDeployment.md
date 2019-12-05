## 爬虫部署

之前已经在[部署爬虫](../Usage/Spider/Deploy.md)中介绍了，爬虫是自动部署在工作节点上的。下面的示意图展示了Crawlab爬虫部署的架构。

![](https://crawlab.oss-cn-hangzhou.aliyuncs.com/v0.3.0/node-deployment.png)

如上图所示，整个爬虫自动部署的生命周期如下(源码在`services/spider.go#InitSpiderService`)：

1. 主节点每5秒，会从爬虫的目录获取爬虫信息，然后更新到数据库（这个过程不涉及文件上传）；
2. 主节点每60秒,从数据库获取所有的爬虫信息，然后将爬虫打包成zip文件，并上传到MongoDB GridFS，并且在MongoDB的`spiders`表里写入`file_id`文件ID；
3. 主节点通过Redis `PubSub`发布消息（`file.upload`事件，包含文件ID）给工作节点，通知工作节点获取爬虫文件；
4. 工作节点接收到获取爬虫文件的消息，从MongoDB GridFS获取zip文件，并解压储存在本地。

这样，所有爬虫将被周期性的部署在工作节点上。

### MongoDB GridFS

GridFS是MongoDB储存大文件（大于16Mb）的文件系统。Crawlab利用GridFS作为了爬虫文件储存的中间媒介，可以让工作节点主动去获取并部署在本地。这样绕开了其他传统传输方式，例如RPC、消息队列、HTTP，因为这几种都要求更复杂也更麻烦的配置和处理。

Crawlab在GridFS上储存文件，会生成两个collection，`files.files`和`files.fs`。前者储存文件的元信息，后者储存文件内容。`spiders`里的`file_id`是指向`files.files`的`_id`。

参考: https://docs.mongodb.com/manual/core/gridfs/
