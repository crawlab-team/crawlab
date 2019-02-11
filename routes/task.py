import json

from celery.utils.log import get_logger
from flask import jsonify
from flask_restful import reqparse, Resource

from app import api
from db.manager import db_manager

logger = get_logger('tasks')
parser = reqparse.RequestParser()

# collection name
COL_NAME = 'tasks'


class TaskApi(Resource):
    col_name = COL_NAME
    parser = reqparse.RequestParser()

    def __init__(self):
        super(TaskApi).__init__()
        self.parser.add_argument('page')
        self.parser.add_argument('page_size')
        self.parser.add_argument('filter')

    def get(self, id=None):
        args = self.parser.parse_args()

        # get item by id
        if id is None:
            # filter
            cond = {}
            if args.get('filter') is not None:
                cond = json.loads(args.filter)

            # page number
            page = 0
            if args.get('page') is not None:
                page = int(args.page)
            else:
                print(args)

            # page size
            page_size = 10
            if args.get('page_size') is not None:
                page = int(args.page_size)

            # total count
            total_count = db_manager.count(col_name=self.col_name, cond=cond)

            # items
            items = db_manager.list(col_name=self.col_name,
                                    cond=cond,
                                    skip=page * page_size,
                                    limit=page_size)
            return jsonify({
                'status': 'ok',
                'total_count': total_count,
                'page': page,
                'page_size': page_size,
                'items': items
            })

        # list items
        else:
            return jsonify(db_manager.get(col_name=self.col_name, id=id))

    def update(self, id=None):
        pass

    def remove(self, id=None):
        pass


# api.add_resource(TaskApi, '/api/task/:id')
api.add_resource(TaskApi,
                 '/api/tasks',
                 '/api/task/:id'
                 )
