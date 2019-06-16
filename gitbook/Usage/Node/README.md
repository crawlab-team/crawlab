## 节点

节点其实就是Celery中的Worker。一个节点运行时会连接到一个任务队列（例如Redis）来接收和运行任务。所有爬虫需要在运行时被部署到节点上，用户在部署前需要定义节点的IP地址和端口（默认为`localhost:8000`）。

1. [查看节点](/Usage/Node/View.md)
2. [修改节点信息](/Usage/Node/Edit.md)
