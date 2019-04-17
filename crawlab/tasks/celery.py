from celery import Celery
# from redisbeat.scheduler import RedisScheduler
from utils.redisbeat import RedisScheduler

# celery app instance
celery_app = Celery(__name__)
celery_app.config_from_object('config')

# RedisBeat scheduler
celery_scheduler = RedisScheduler(app=celery_app)
