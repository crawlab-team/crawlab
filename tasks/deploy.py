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
def deploy_spider(id):
    pass
