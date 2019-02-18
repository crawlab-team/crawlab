import os
import sys
import threading

from celery import Celery

app = Celery(__name__)
app.config_from_object('config.celery')

import tasks.spider
import tasks.deploy

if __name__ == '__main__':
    if sys.platform == 'windows':
        app.start(argv=['tasks', 'worker', '-P', 'eventlet', '-E', '-l', 'INFO'])
    else:
        app.start(argv=['tasks', 'worker', '-E', '-l', 'INFO'])
