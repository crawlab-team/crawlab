import json
# from celery.utils.log import get_logger
import os
import shutil

from bson import ObjectId
from flask_restful import reqparse, Resource

from app import api
from config import PROJECT_FILE_FOLDER
from db.manager import db_manager
from routes.base import BaseApi
from tasks.spider import execute_spider
from utils import jsonify


class SpiderApi(BaseApi):
    col_name = 'spiders'

    arguments = (
        ('name', str),
        ('cmd', str),
        ('src', str),
        ('type', int),
        ('lang', int),
    )

    def crawl(self, id):
        job = execute_spider.delay(id)
        # print('crawl: %s' % id)
        return {
            'code': 200,
            'status': 'ok',
            'task': {
                'id': job.id,
                'status': job.status
            }
        }

    def deploy(self, id):
        # get spider given the id
        spider = db_manager.get(col_name=self.col_name, id=id)
        if spider is None:
            return

        # get latest version
        latest_version = db_manager.get_latest_version(spider_id=id)

        # initialize version if no version found
        if latest_version is None:
            latest_version = 0

        # make source / destination
        src = spider.get('src')
        dst = os.path.join(PROJECT_FILE_FOLDER, str(spider.get('_id')), str(latest_version + 1))

        # copy files
        try:
            shutil.copytree(src=src, dst=dst)
            return {
                'code': 200,
                'status': 'ok',
                'message': 'deploy success'
            }
        except Exception as err:
            print(err)
            return {
                'code': 500,
                'status': 'ok',
                'error': str(err)
            }
        finally:
            version = latest_version + 1
            db_manager.save('deploys', {
                'spider_id': ObjectId(id),
                'version': version,
                'node_id': None  # TODO: deploy to corresponding node
            })


api.add_resource(SpiderApi,
                 '/api/spiders',
                 '/api/spiders/<string:id>',
                 '/api/spiders/<string:id>/<string:action>')
