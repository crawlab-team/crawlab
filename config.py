# project variables
PROJECT_SOURCE_FILE_FOLDER = '/Users/yeqing/projects/crawlab/spiders'
PROJECT_DEPLOY_FILE_FOLDER = '/var/crawlab'
PROJECT_LOGS_FOLDER = '/var/logs/crawlab'
PROJECT_TMP_FOLDER = '/tmp'

# celery variables
# BROKER_URL = 'redis://localhost:6379/0'
# CELERY_RESULT_BACKEND = 'mongodb://localhost:27017/'
BROKER_URL = 'redis://192.168.99.100:6379/0'
CELERY_RESULT_BACKEND = 'mongodb://192.168.99.100:27017/'
CELERY_MONGODB_BACKEND_SETTINGS = {
    'database': 'crawlab_test',
    'taskmeta_collection': 'tasks_celery',
}
CELERY_TIMEZONE = 'Asia/Shanghai'

# flower variables
FLOWER_API_ENDPOINT = 'http://localhost:5555/api'

# database variables
# MONGO_HOST = 'localhost'
MONGO_HOST = '192.168.99.100'
MONGO_PORT = 27017
# MONGO_USER = 'test'
# MONGO_PASS = 'test'
MONGO_DB = 'crawlab_test'

# flask variables
DEBUG = True
FLASK_HOST = '0.0.0.0'
FLASK_PORT = 5000
# SERVER_NAME = '0.0.0.0:5000'
