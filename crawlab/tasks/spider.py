import os
from datetime import datetime

from bson import ObjectId
from celery.utils.log import get_logger

from config import PROJECT_DEPLOY_FILE_FOLDER, PROJECT_LOGS_FOLDER
from constants.task import TaskStatus
from db.manager import db_manager
from .celery import celery_app
import subprocess

logger = get_logger(__name__)


@celery_app.task(bind=True)
def execute_spider(self, id: str):
    task_id = self.request.id
    hostname = self.request.hostname
    spider = db_manager.get('spiders', id=id)
    command = spider.get('cmd')

    current_working_directory = os.path.join(PROJECT_DEPLOY_FILE_FOLDER, str(spider.get('_id')))

    # log info
    logger.info('task_id: %s' % task_id)
    logger.info('hostname: %s' % hostname)
    logger.info('current_working_directory: %s' % current_working_directory)
    logger.info('spider_id: %s' % id)
    logger.info(command)

    # make sure the log folder exists
    log_path = os.path.join(PROJECT_LOGS_FOLDER, id)
    if not os.path.exists(log_path):
        os.makedirs(log_path)

    # open log file streams
    log_file_path = os.path.join(log_path, '%s.log' % datetime.now().strftime('%Y%m%d%H%M%S'))
    stdout = open(log_file_path, 'a')
    stderr = open(log_file_path, 'a')

    # create a new task
    db_manager.update_one('tasks', id=task_id, values={
        'start_ts': datetime.utcnow(),
        'node_id': hostname,
        'hostname': hostname,
        'log_file_path': log_file_path,
        'status': TaskStatus.STARTED
    })

    # start the process and pass params as env variables
    env = os.environ.copy()
    env['CRAWLAB_TASK_ID'] = task_id
    if spider.get('col'):
        env['CRAWLAB_COLLECTION'] = spider.get('col')
    p = subprocess.Popen(command.split(' '),
                         stdout=stdout.fileno(),
                         stderr=stderr.fileno(),
                         cwd=current_working_directory,
                         env=env,
                         bufsize=1)

    # get output from the process
    _stdout, _stderr = p.communicate()

    # get return code
    code = p.poll()
    if code == 0:
        status = TaskStatus.SUCCESS
    else:
        status = TaskStatus.FAILURE

    # save task when the task is finished
    db_manager.update_one('tasks', id=task_id, values={
        'start_ts': datetime.utcnow(),
        'node_id': hostname,
        'hostname': hostname,
        'log_file_path': log_file_path,
        'finish_ts': datetime.utcnow(),
        'status': status
    })
    task = db_manager.get('tasks', id=id)

    # close log file streams
    stdout.flush()
    stderr.flush()
    stdout.close()
    stderr.close()

    return task
