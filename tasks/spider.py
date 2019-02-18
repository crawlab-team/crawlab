import os
import sys
from datetime import datetime

import requests
from celery.utils.log import get_logger

from config import PROJECT_FILE_FOLDER, PROJECT_LOGS_FOLDER
from db.manager import db_manager
from tasks import app
import subprocess

logger = get_logger(__name__)


@app.task
def execute_spider(id: str):
    spider = db_manager.get('spiders', id=id)
    latest_version = db_manager.get_latest_version(spider_id=id)
    command = spider.get('cmd')
    current_working_directory = os.path.join(PROJECT_FILE_FOLDER, str(spider.get('_id')), str(latest_version))

    # log info
    logger.info('spider_id: %s' % id)
    logger.info('version: %s' % latest_version)
    logger.info(command)

    # make sure the log folder exists
    log_path = os.path.join(PROJECT_LOGS_FOLDER, id, str(latest_version))
    if not os.path.exists(log_path):
        os.makedirs(log_path)

    # execute the command
    p = subprocess.Popen(command,
                         shell=True,
                         stdout=subprocess.PIPE,
                         stderr=subprocess.PIPE,
                         cwd=current_working_directory,
                         bufsize=1)

    # output the log file
    log_file_path = os.path.join(log_path, '%s.txt' % datetime.now().strftime('%Y%m%d%H%M%S'))
    with open(log_file_path, 'a') as f:
        for line in p.stdout.readlines():
            f.write(line.decode('utf-8') + '\n')


@app.task
def get_baidu_html(keyword: str):
    res = requests.get('http://www.baidu.com')
    return res.content.decode('utf-8')
