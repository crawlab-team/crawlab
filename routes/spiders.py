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


class SpiderExecutorApi(Resource):
    col_name = 'spiders'

    def post(self, id):
        args = parser.parse_args()
        job = execute_spider.delay(args.spider_name)
        return {
            'id': job.id,
            'status': job.status,
            'spider_name': args.spider_name,
            'result': job.get(timeout=5)
        }


api.add_resource(SpiderExecutorApi, '/api/spiders/<string:id>/crawl')
api.add_resource(SpiderApi,
                 '/api/spiders',
                 '/api/spiders/<string:id>')
