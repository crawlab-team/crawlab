import os
import sys
from datetime import datetime

import requests
from celery.utils.log import get_logger

from config import PROJECT_DEPLOY_FILE_FOLDER, PROJECT_LOGS_FOLDER
from db.manager import db_manager
from .celery import celery_app
import subprocess

logger = get_logger(__name__)


@celery_app.task
def deploy_spider(id):
    pass
