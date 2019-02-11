from celery.utils.log import get_logger
from flask_restful import reqparse, Resource
from tasks.spider import execute_spider

logger = get_logger('tasks')
parser = reqparse.RequestParser()
parser.add_argument('spider_name', type=str)


class SpiderApi(Resource):
    pass


class SpiderExecutorApi(Resource):
    def get(self):
        args = parser.parse_args()
        job = execute_spider.delay(args.spider_name)
        return {
            'id': job.id,
            'status': job.status,
            'spider_name': args.spider_name,
            'result': job.get(timeout=5)
        }
