import os
import sys

import requests
from celery.utils.log import get_logger

from config import PROJECT_FILE_FOLDER
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
    p = subprocess.Popen(command,
                         shell=True,
                         stdout=subprocess.PIPE,
                         stderr=subprocess.STDOUT,
                         cwd=current_working_directory,
                         bufsize=1)
    for i in iter(p.stdout.readline, 'b'):
        yield i


@app.task
def get_baidu_html(keyword: str):
    res = requests.get('http://www.baidu.com')
    return res.content.decode('utf-8')
