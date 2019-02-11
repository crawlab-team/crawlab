import json
# from celery.utils.log import get_logger
from flask_restful import reqparse, Resource

from app import api
from db.manager import db_manager
from tasks.spider import execute_spider

# logger = get_logger('tasks')
parser = reqparse.RequestParser()
parser.add_argument('spider_name', type=str)

# collection name
COL_NAME = 'spiders'


class SpiderApi(Resource):
    col_name = COL_NAME

    def get(self, id=None):
        args = parser.parse_args()
        cond = {}
        if args.filter is not None:
            cond = json.loads(args.filter)
        if id is None:
            return db_manager.list(col_name=self.col_name, cond=cond, page=args.page, page_size=args.page_size)
        else:
            return db_manager.get(col_name=self.col_name, id=id)

    def list(self):
        args = parser.parse_args()
        cond = {}
        if args.filter is not None:
            cond = json.loads(args.filter)
        return db_manager.list(col_name=self.col_name, cond=cond, page=args.page, page_size=args.page_size)

    def update(self, id=None):
        pass

    def remove(self, id=None):
        pass


class SpiderExecutorApi(Resource):
    col_name = COL_NAME

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
