import json
# from celery.utils.log import get_logger
from flask_restful import reqparse, Resource

from app import api
from db.manager import db_manager
from routes.base import BaseApi
from tasks.spider import execute_spider

# logger = get_logger('tasks')
parser = reqparse.RequestParser()
parser.add_argument('spider_name', type=str)


class SpiderApi(BaseApi):
    col_name = 'spiders'

    arguments = (
        ('spider_name', str),
        ('spider_type', int),
        ('lang_type', int),
        ('execute_cmd', str),
        ('src_file_path', str),
    )

    def crawl(self, id):
        print('crawl: %s' % id)

    def deploy(self, id):
        print('deploy: %s' % id)


api.add_resource(SpiderApi,
                 '/api/spiders',
                 '/api/spiders/<string:id>',
                 '/api/spiders/<string:id>/<string:action>')
