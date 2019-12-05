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
npm install -g yarn
cd frontend
yarn install
```

安装后端所需库。

```bash
cd ../backend
go install ./...
```

### 配置

修改配置文件`./backend/config.yaml`。配置文件是以`yaml`的格式。配置详情请见[配置Crawlab](../Config/README.md)。

### 构建前端

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
    root    /path/to/dist;
    index    index.html;
}
```

其中，`root`是静态文件的根目录，这里是`npm`打包好后的静态文件。

现在，只需要启动`nginx`服务就完成了启动前端服务。

```bash
nginx reload
```

### 构建后端

执行以下命令。

```bash
cd ../backend
go build
```

`go build`命令会将Golang代码打包为一个执行文件，默认在`$GOPATH/bin`里。

### 启动服务

这里是指启动后端服务。执行以下命令。

```bash
$GOPATH/bin/crawlab
```

然后在浏览器中输入`http://localhost:8080`就可以看到界面了。
