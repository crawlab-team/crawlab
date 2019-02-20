BROKER_URL = 'redis://localhost:6379/0'
# BROKER_URL = 'mongodb://localhost:27017/'
CELERY_RESULT_BACKEND = 'mongodb://localhost:27017/'
# CELERY_RESULT_BACKEND = 'redis://localhost:6379/1'
# CELERY_TASK_SERIALIZER = 'json'
# CELERY_RESULT_SERIALIZER = 'json'
# CELERY_TASK_RESULT_EXPIRES = 60 * 60 * 24  # 任务过期时间
CELERY_MONGODB_BACKEND_SETTINGS = {
    'database': 'crawlab_test',
    'taskmeta_collection': 'tasks_celery',
}

FLOWER_API_ENDPOINT = 'http://localhost:5555/api'
