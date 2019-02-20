import sys
from celery import Celery

from app import celery_app

# import necessary tasks
import tasks.spider
import tasks.deploy

if __name__ == '__main__':
    if sys.platform == 'windows':
        celery_app.start(argv=['tasks', 'worker', '-P', 'eventlet', '-E', '-l', 'INFO'])
    else:
        celery_app.start(argv=['tasks', 'worker', '-E', '-l', 'INFO'])
