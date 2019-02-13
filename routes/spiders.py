import json
# from celery.utils.log import get_logger
import os
import shutil

from flask_restful import reqparse, Resource

from app import api
from config import PROJECT_FILE_FOLDER
from db.manager import db_manager
from routes.base import BaseApi
from tasks.spider import execute_spider


class SpiderApi(BaseApi):
    col_name = 'spiders'

    arguments = (
        ('spider_name', str),
        ('cmd', str),
        ('src', str),
        ('spider_type', int),
        ('lang_type', int),
    )

    def crawl(self, id):
        print('crawl: %s' % id)

    def deploy(self, id):
        args = self.parser.parse_args()
        spider = db_manager.get(col_name=self.col_name, id=id)
        latest_version = db_manager.get_latest_version(id=id)
        src = args.get('src')
        dst = os.path.join(PROJECT_FILE_FOLDER, str(spider._id), latest_version + 1)
        if not os.path.exists(dst):
            os.mkdir(dst)
        shutil.copytree(src=src, dst=dst)


api.add_resource(SpiderApi,
                 '/api/spiders',
                 '/api/spiders/<string:id>',
                 '/api/spiders/<string:id>/<string:action>')
