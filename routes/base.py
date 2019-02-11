from flask_restful import reqparse, Resource

from db.manager import db_manager
from utils import jsonify


class BaseApi(Resource):
    col_name = 'tmp'
    parser = reqparse.RequestParser()
    arguments = []

    def __init__(self):
        super(BaseApi).__init__()
        self.parser.add_argument('page', type=int)
        self.parser.add_argument('page_size', type=int)
        self.parser.add_argument('filter', type=dict)

    def get(self, id=None):
        args = self.parser.parse_args()

        # get item by id
        if id is None:
            # filter
            cond = {}
            if args.get('filter') is not None:
                cond = args.filter
                # cond = json.loads(args.filter)

            # page number
            page = 1
            if args.get('page') is not None:
                page = args.page
                # page = int(args.page)

            # page size
            page_size = 10
            if args.get('page_size') is not None:
                page_size = args.page_size
                # page = int(args.page_size)

            # TODO: sort functionality

            # total count
            total_count = db_manager.count(col_name=self.col_name, cond=cond)

            # items
            items = db_manager.list(col_name=self.col_name,
                                    cond=cond,
                                    skip=(page - 1) * page_size,
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
