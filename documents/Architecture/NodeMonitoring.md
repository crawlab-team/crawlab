## 节点监控

Crawlab的节点监控是通过Redis来完成的。原理如下图。

![](https://crawlab.oss-cn-hangzhou.aliyuncs.com/v0.3.0/node-monitoring.png)

工作节点会不断更新心跳信息在Redis上，利用`HSET nodes <node_id> <msg>`，心跳信息`<msg>`包含节点MAC地址，IP地址，当前时间戳，

主节点会周期性获取Redis上的工作节点心跳信息。如果有工作节点的时间戳在60秒之前，则考虑该节点为离线状态，会在Redis中删除该节点的信息，并在MongoDB中设置为"离线"；如果时间戳在过去60秒之内，则保留该节点信息，在MongoDB中设置为"在线"。
