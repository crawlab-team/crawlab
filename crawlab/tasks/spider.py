import os
import sys
from datetime import datetime
from time import sleep

from bson import ObjectId
from pymongo import ASCENDING, DESCENDING

from config import PROJECT_DEPLOY_FILE_FOLDER, PROJECT_LOGS_FOLDER, PYTHON_ENV_PATH, MONGO_HOST, MONGO_PORT, MONGO_DB
from constants.task import TaskStatus
from db.manager import db_manager
from .celery import celery_app
import subprocess
from utils.log import other as logger

BASE_DIR = os.path.abspath(os.path.join(os.path.dirname(__file__), '..'))


def get_task(id: str):
    i = 0
    while i < 5:
        task = db_manager.get('tasks', id=id)
        if task is not None:
            return task
        i += 1
        sleep(1)
    return None


@celery_app.task(bind=True)
def execute_spider(self, id: str, params: str = None):
    """
    Execute spider task.
    :param self:
    :param id: task_id
    """
    task_id = self.request.id
    hostname = self.request.hostname
    spider = db_manager.get('spiders', id=id)
    command = spider.get('cmd')

    # if start with python, then use sys.executable to execute in the virtualenv
    if command.startswith('python '):
        command = command.replace('python ', sys.executable + ' ')

    # if start with scrapy, then use sys.executable to execute scrapy as module in the virtualenv
    elif command.startswith('scrapy '):
        command = command.replace('scrapy ', sys.executable + ' -m scrapy ')

    # pass params to the command
    if params is not None:
        command += ' ' + params

    # get task object and return if not found
    task = get_task(task_id)
    if task is None:
        return

    # current working directory
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

    # update task status as started
    db_manager.update_one('tasks', id=task_id, values={
        'start_ts': datetime.utcnow(),
        'node_id': hostname,
        'hostname': hostname,
        'log_file_path': log_file_path,
        'status': TaskStatus.STARTED
    })

    # pass params as env variables
    env = os.environ.copy()

    # custom environment variables
    if spider.get('envs'):
        for _env in spider.get('envs'):
            env[_env['name']] = _env['value']

    # task id environment variable
    env['CRAWLAB_TASK_ID'] = task_id

    # collection environment variable
    if spider.get('col'):
        env['CRAWLAB_COLLECTION'] = spider.get('col')

        # create index to speed results data retrieval
        db_manager.create_index(spider.get('col'), [('task_id', ASCENDING)])

    # start process
    cmd_arr = command.split(' ')
    cmd_arr = list(filter(lambda x: x != '', cmd_arr))
    try:
        p = subprocess.Popen(cmd_arr,
                             stdout=stdout.fileno(),
                             stderr=stderr.fileno(),
                             cwd=current_working_directory,
                             env=env,
                             bufsize=1)

        # update pid
        db_manager.update_one(col_name='tasks', id=task_id, values={
            'pid': p.pid
        })

        # get output from the process
        _stdout, _stderr = p.communicate()

        # get return code
        code = p.poll()
        if code == 0:
            status = TaskStatus.SUCCESS
        else:
            status = TaskStatus.FAILURE
    except Exception as err:
        logger.error(err)
        stderr.write(str(err))
        status = TaskStatus.FAILURE

    # save task when the task is finished
    finish_ts = datetime.utcnow()
    db_manager.update_one('tasks', id=task_id, values={
        'finish_ts': finish_ts,
        'duration': (finish_ts - task['create_ts']).total_seconds(),
        'status': status
    })
    task = db_manager.get('tasks', id=id)

    # close log file streams
    stdout.flush()
    stderr.flush()
    stdout.close()
    stderr.close()

    return task


@celery_app.task(bind=True)
def execute_config_spider(self, id: str, params: str = None):
    task_id = self.request.id
    hostname = self.request.hostname
    spider = db_manager.get('spiders', id=id)

    # get task object and return if not found
    task = get_task(task_id)
    if task is None:
        return

    # current working directory
    current_working_directory = os.path.join(BASE_DIR, 'spiders')

    # log info
    logger.info('task_id: %s' % task_id)
    logger.info('hostname: %s' % hostname)
    logger.info('current_working_directory: %s' % current_working_directory)
    logger.info('spider_id: %s' % id)

    # make sure the log folder exists
    log_path = os.path.join(PROJECT_LOGS_FOLDER, id)
    if not os.path.exists(log_path):
        os.makedirs(log_path)

    # open log file streams
    log_file_path = os.path.join(log_path, '%s.log' % datetime.now().strftime('%Y%m%d%H%M%S'))
    stdout = open(log_file_path, 'a')
    stderr = open(log_file_path, 'a')

    # update task status as started
    db_manager.update_one('tasks', id=task_id, values={
        'start_ts': datetime.utcnow(),
        'node_id': hostname,
        'hostname': hostname,
        'log_file_path': log_file_path,
        'status': TaskStatus.STARTED
    })

    # pass params as env variables
    env = os.environ.copy()

    # custom environment variables
    if spider.get('envs'):
        for _env in spider.get('envs'):
            env[_env['name']] = _env['value']

    # task id environment variable
    env['CRAWLAB_TASK_ID'] = task_id

    # collection environment variable
    if spider.get('col'):
        env['CRAWLAB_COLLECTION'] = spider.get('col')

        # create index to speed results data retrieval
        db_manager.create_index(spider.get('col'), [('task_id', ASCENDING)])

    # mongodb environment variables
    env['MONGO_HOST'] = MONGO_HOST
    env['MONGO_PORT'] = str(MONGO_PORT)
    env['MONGO_DB'] = MONGO_DB

    cmd_arr = [
        sys.executable,
        '-m',
        'scrapy',
        'crawl',
        'config_spider'
    ]
    try:
        p = subprocess.Popen(cmd_arr,
                             stdout=stdout.fileno(),
                             stderr=stderr.fileno(),
                             cwd=current_working_directory,
                             env=env,
                             bufsize=1)

        # update pid
        db_manager.update_one(col_name='tasks', id=task_id, values={
            'pid': p.pid
        })

        # get output from the process
        _stdout, _stderr = p.communicate()

        # get return code
        code = p.poll()
        if code == 0:
            status = TaskStatus.SUCCESS
        else:
            status = TaskStatus.FAILURE
    except Exception as err:
        logger.error(err)
        stderr.write(str(err))
        status = TaskStatus.FAILURE

    # save task when the task is finished
    finish_ts = datetime.utcnow()
    db_manager.update_one('tasks', id=task_id, values={
        'finish_ts': finish_ts,
        'duration': (finish_ts - task['create_ts']).total_seconds(),
        'status': status
    })
    task = db_manager.get('tasks', id=id)

    # close log file streams
    stdout.flush()
    stderr.flush()
    stdout.close()
    stderr.close()

    return task
