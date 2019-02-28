# project variables
PROJECT_SOURCE_FILE_FOLDER = '/Users/yeqing/projects/crawlab/spiders'
PROJECT_DEPLOY_FILE_FOLDER = '/var/crawlab'
PROJECT_LOGS_FOLDER = '/Users/yeqing/projects/crawlab/logs/crawlab'
PROJECT_TMP_FOLDER = '/tmp'

# celery variables
BROKER_URL = 'redis://localhost:6379/0'
CELERY_RESULT_BACKEND = 'mongodb://localhost:27017/'
CELERY_MONGODB_BACKEND_SETTINGS = {
    'database': 'crawlab_test',
    'taskmeta_collection': 'tasks_celery',
}
FLOWER_API_ENDPOINT = 'http://localhost:5555/api'

# database variables
MONGO_HOST = 'localhost'
MONGO_PORT = 27017
# MONGO_USER = 'test'
# MONGO_PASS = 'test'
MONGO_DB = 'crawlab_test'

# flask variables
DEBUG = True
