from celery import Celery

app = Celery(__name__)
app.config_from_object('config.celery')

import tasks.spider

if __name__ == '__main__':
    app.start(argv=['tasks.spider', 'worker', '-P', 'eventlet', '-E', '-l', 'INFO'])
