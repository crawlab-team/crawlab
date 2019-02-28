import os
import sys
from datetime import datetime

import requests
from bson import ObjectId
from celery.utils.log import get_logger

from config import PROJECT_DEPLOY_FILE_FOLDER, PROJECT_LOGS_FOLDER
from db.manager import db_manager
from .celery import celery_app
import subprocess

logger = get_logger(__name__)


@celery_app.task(bind=True)
def execute_spider(self, id: str, node_id: str):
    task_id = self.request.id
    hostname = self.request.hostname
    spider = db_manager.get('spiders', id=id)
    latest_version = db_manager.get_latest_version(spider_id=id)
    command = spider.get('cmd')
    current_working_directory = os.path.join(PROJECT_DEPLOY_FILE_FOLDER, str(spider.get('_id')), str(latest_version))

    # log info
    logger.info('spider_id: %s' % id)
    logger.info('version: %s' % latest_version)
    logger.info(command)

    # make sure the log folder exists
    log_path = os.path.join(PROJECT_LOGS_FOLDER, id, str(latest_version))
    if not os.path.exists(log_path):
        os.makedirs(log_path)

    # open log file streams
    log_file_path = os.path.join(log_path, '%s.log' % datetime.now().strftime('%Y%m%d%H%M%S'))
    stdout = open(log_file_path, 'a')
    stderr = open(log_file_path, 'a')

    # create a new task
    db_manager.save('tasks', {
        '_id': task_id,
        'spider_id': ObjectId(id),
        'create_ts': datetime.now(),
        'node_id': node_id,
        'hostname': hostname,
        'log_file_path': log_file_path,
        'spider_version': latest_version
    })

    # execute the command
    p = subprocess.Popen(command.split(' '),
                         stdout=stdout.fileno(),
                         stderr=stderr.fileno(),
                         cwd=current_working_directory,
                         bufsize=1)

    # get output from the process
    _stdout, _stderr = p.communicate()

    # save task when the task is finished
    db_manager.update_one('tasks', id=task_id, values={
        'finish_ts': datetime.now(),
    })
    task = db_manager.get('tasks', id=id)

    # close log file streams
    stdout.flush()
    stderr.flush()
    stdout.close()
    stderr.close()

    return task
