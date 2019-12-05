## 节点通信

这里的通信主要是指节点间的即时通信，即没有显著的延迟（[爬虫部署](./SpiderDeployment.md)和[任务执行](./TaskExecution.md)是通过轮训来完成的，不在此列）。

通信主要由Redis来完成。以下为节点通信原理示意图。

![](https://crawlab.oss-cn-hangzhou.aliyuncs.com/v0.3.0/node-communication.png)

各个节点会通过Redis的`PubSub`功能来做相互通信。

所谓`PubSub`，简单来说是一个发布订阅模式。订阅者（Subscriber）会在Redis上订阅（Subscribe）一个通道，其他任何一个节点都可以作为发布者（Publisher）在该通道上发布（Publish）消息。

在Crawlab中，主节点会订阅`nodes:master`通道，其他节点如果需要向主节点发送消息，只需要向`nodes:master`发布消息就可以了。同理，各工作节点会各自订阅一个属于自己的通道`nodes:<node_id>`（`node_id`是MongoDB里的节点ID，是MongoDB ObjectId），如果需要给工作节点发送消息，只需要发布消息到该通道就可以了。

一个网络请求的简单过程如下:
1. 客户端（前端应用）发送请求给主节点（API）；
2. 主节点通过Redis `PubSub`的`<nodes:<node_id>`通道发布消息给相应的工作节点；
3. 工作节点收到消息之后，执行一些操作，并将相应的消息通过`<nodes:master>`通道发布给主节点；
4. 主节点收到消息之后，将消息返回给客户端。

不是所有节点通信都是双向的，也就是说，主节点只会单方面对工作节点通信，工作节点并不会返回响应给主节点，所谓的单向通信。以下是Crawlab的通信类型。

操作名称 | 通信类别
--- | ---
获取日志 | 双向通信
获取系统信息 | 双向通信
取消任务 | 单向通信
通知工作节点向GridFS获取爬虫文件 | 单向通信

### `chan`和`goroutine`

如果您在阅读Crawlab源码，会发现节点通信中有大量的`chan`语法，这是Golang的一个并发特性。

`chan`表示为一个通道，在Golang中分为无缓冲和有缓冲的通道，我们用了无缓冲通道来阻塞协程，只有当`chan`接收到信号（`chan <- "some signal"`），该阻塞才会释放，协程进行下一步操作）。在请求响应模式中，如果为双向通信，主节点收到请求后会起生成一个无缓冲通道来阻塞该请求，当收到来自工作节点的消息后，向该无缓冲通道赋值，阻塞释放，返回响应给客户端。

`go`命令会起一个`goroutine`（协程）来完成并发，配合`chan`，该协程可以利用无缓冲通道挂起，等待信号执行接下来的操作。任务取消就是`go`+`chan`来实现的。有兴趣的读者可以参考一下[源码](https://github.com/tikazyq/crawlab/blob/master/backend/services/task.go#L136)。

### Redis PubSub

这是Redis版发布／订阅消息模式的一种实现。其用法非常简单：
1. 订阅者利用`SUBSCRIBE channel1 channel2 ...`来订阅一个或多个频道；
2. 发布者利用`PUBLISH channelx message`来发布消息给该频道的订阅者。

Redis的`PubSub`可以用作广播模式，即一个发布者对应多个订阅者。而在Crawlab中，我们只有一个订阅者对应一个发布者的情况（主节点->工作节点：`nodes:<node_id>`）或一个订阅者对应多个发布者的情况（工作节点->主节点：`nodes:master>`）。这是为了方便双向通信。

参考：https://redis.io/topics/pubsub
