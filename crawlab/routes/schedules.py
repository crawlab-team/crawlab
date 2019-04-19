import json

import requests

from constants.task import TaskStatus
from db.manager import db_manager
from routes.base import BaseApi
from tasks.scheduler import scheduler
from utils import jsonify
from utils.spider import get_spider_col_fields


class ScheduleApi(BaseApi):
    col_name = 'schedules'

    arguments = (
        ('name', str),
        ('description', str),
        ('cron', str),
        ('spider_id', str)
    )

    def after_update(self, id: str = None):
        scheduler.update()
