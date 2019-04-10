import os
from datetime import datetime

from bson import ObjectId
from config import PROJECT_DEPLOY_FILE_FOLDER, PROJECT_LOGS_FOLDER, PYTHON_ENV_PATH
from constants.task import TaskStatus
from db.manager import db_manager
from .celery import celery_app
import subprocess
from utils.log import other as logger


@celery_app.task(bind=True)
def execute_spider(self, id: str):
    """
    Execute spider task.
    :param self:
    :param id: task_id
    """
    task_id = self.request.id
    hostname = self.request.hostname
    spider = db_manager.get('spiders', id=id)
    command = spider.get('cmd')
    if command.startswith("env"):
        command = PYTHON_ENV_PATH + command.replace("env", "")

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

    # start process
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
