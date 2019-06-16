## 直接部署

直接部署是之前没有Docker时的部署方式，相对于Docker部署来说有些繁琐。但了解如何直接部署可以帮助更深入地理解Docker是如何构建Crawlab镜像的。这里简单介绍一下。

### 拉取代码

首先是将github上的代码拉取到本地。

```bash
git clone https://github.com/tikazyq/crawlab
```

### 安装

安装前端所需库。

```bash
npm install -g yarn pm2
cd frontend
yarn install
```

安装后端所需库。

```bash
cd ../crawlab
pip install -r requirements
```

### 配置

分别配置前端配置文件`./frontend/.env.production`和后端配置文件`./crawlab/config/config.py`。分别需要对部署后API地址以及数据库地址进行配置。

### 构建

这里的构建是指前端构建，需要执行以下命令。

```bash
cd ../frontend
npm run build:prod
```

构建完成后，会在`./frontend`目录下创建一个`dist`文件夹，里面是打包好后的静态文件。

### Nginx

安装`nginx`，在`ubuntu 16.04`是以下命令。

```bash
sudo apt-get install nginx
```

添加`/etc/nginx/conf.d/crawlab.conf`文件，输入以下内容。

```
server {
    listen    8080;
    server_name    dev.crawlab.com;
    root    /home/yeqing/jenkins_home/workspace/crawlab_develop/frontend/dist;
    index    index.html;
}
```

其中，`root`是静态文件的根目录，这里是`npm`打包好后的静态文件。

现在，只需要启动`nginx`服务就完成了启动前端服务。

```bash
nginx reload
```

### 启动服务

这里是指启动后端服务。我们用`pm2`来管理进程。执行以下命令。

```bash
pm2 start app.py # API服务
pm2 start worker.py # Worker
pm2 start flower.py # Flower
```

这样，`pm2`会启动3个守护进程来管理这3个服务。我们如果想看后端服务的日志的话，可以执行以下命令。

```bash
pm2 logs [app]
```

然后在浏览器中输入`http://localhost:8080`就可以看到界面了。