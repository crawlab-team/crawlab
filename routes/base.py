import json

from celery.utils.log import get_logger
from flask_restful import reqparse, Resource

from app import api
from db.manager import db_manager

logger = get_logger('tasks')
parser = reqparse.RequestParser()
parser.add_argument('task_name', type=str)

# collection name
COL_NAME = 'test'


class BaseApi(Resource):
    col_name = 'base'

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

