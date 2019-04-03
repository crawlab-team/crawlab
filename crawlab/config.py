# project variables
# 爬虫源码路径
PROJECT_SOURCE_FILE_FOLDER = '../spiders'
# 配置python虚拟环境的路径
PYTHON_ENV_PATH="/Users/chennan/Desktop/2019/env/bin/python"
# 爬虫部署路径
PROJECT_DEPLOY_FILE_FOLDER = '../deployfile'

PROJECT_LOGS_FOLDER = '../deployfile/logs'
PROJECT_TMP_FOLDER = '/tmp'

# celery variables
BROKER_URL = 'redis://127.0.0.1:6379/0'
CELERY_RESULT_BACKEND = 'mongodb://127.0.0.1:27017/'
CELERY_MONGODB_BACKEND_SETTINGS = {
    'database': 'crawlab_test',
    'taskmeta_collection': 'tasks_celery',
}
CELERY_TIMEZONE = 'Asia/Shanghai'
CELERY_ENABLE_UTC = True

# flower variables
FLOWER_API_ENDPOINT = 'http://localhost:5555/api'

# database variables
MONGO_HOST = '127.0.0.1'
MONGO_PORT = 27017
MONGO_DB = 'crawlab_test'

# flask variables
DEBUG = True
FLASK_HOST = '127.0.0.1'
FLASK_PORT = 8000

