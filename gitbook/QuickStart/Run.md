# 运行

在运行之前需要对Crawlab进行一些配置，配置文件为`config.py`。

```python
# project variables
PROJECT_SOURCE_FILE_FOLDER = '/Users/yeqing/projects/crawlab/spiders' # 爬虫源码根目录
PROJECT_DEPLOY_FILE_FOLDER = '/var/crawlab'  # 爬虫部署根目录
PROJECT_LOGS_FOLDER = '/var/logs/crawlab'  # 日志目录
PROJECT_TMP_FOLDER = '/tmp'  # 临时文件目录

# celery variables
BROKER_URL = 'redis://192.168.99.100:6379/0'  # 中间者URL，连接redis
CELERY_RESULT_BACKEND = 'mongodb://192.168.99.100:27017/'  # CELERY后台URL
CELERY_MONGODB_BACKEND_SETTINGS = {
    'database': 'crawlab_test',
    'taskmeta_collection': 'tasks_celery',
}
CELERY_TIMEZONE = 'Asia/Shanghai'
CELERY_ENABLE_UTC = True

# flower variables
FLOWER_API_ENDPOINT = 'http://localhost:5555/api'  # Flower服务地址

# database variables
MONGO_HOST = '192.168.99.100'
MONGO_PORT = 27017
MONGO_DB = 'crawlab_test'

# flask variables
DEBUG = True
FLASK_HOST = '127.0.0.1'
FLASK_PORT = 8000
```

启动后端API，也就是一个Flask App，可以直接启动，或者用gunicorn代替。

```bash
python app.py
```

启动Flower服务（抱歉目前集成Flower到App服务中，必须单独启动来获取节点信息，后面的版本不需要这个操作）。

```bash
python ./bin/run_flower.py
```

启动本地Worker。在其他节点中如果想只是想执行任务的话，只需要启动这一个服务就可以了。

```bash
python ./bin/run_worker.py
```

启动前端服务器。

```bash
cd ../frontend
npm run serve
```
