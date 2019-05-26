import sys
import os

# make sure the working directory is in system path
file_dir = os.path.dirname(os.path.realpath(__file__))
root_path = os.path.abspath(os.path.join(file_dir, '..'))
sys.path.append(root_path)

from tasks.celery import celery_app

# import necessary tasks
import tasks.spider
import tasks.deploy

if __name__ == '__main__':
    if 'win' in sys.platform:
        celery_app.start(argv=['tasks', 'worker', '-P', 'eventlet', '-E', '-l', 'INFO'])
    else:
        celery_app.start(argv=['tasks', 'worker', '-E', '-l', 'INFO'])
