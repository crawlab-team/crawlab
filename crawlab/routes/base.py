from flask_restful import reqparse, Resource
# from flask_restplus import reqparse, Resource

from db.manager import db_manager
from utils import jsonify

DEFAULT_ARGS = [
    'page_num',
    'page_size',
    'filter'
]


class BaseApi(Resource):
    col_name = 'tmp'
    parser = reqparse.RequestParser()
    arguments = []

    def __init__(self):
        super(BaseApi).__init__()
        self.parser.add_argument('page_num', type=int)
        self.parser.add_argument('page_size', type=int)
        self.parser.add_argument('filter', type=dict)

        for arg, type in self.arguments:
            self.parser.add_argument(arg, type=type)

    def get(self, id=None, action=None):
        import pdb
        pdb.set_trace()
        args = self.parser.parse_args()

        # action by id
        if action is not None:
            if not hasattr(self, action):
                return {
                           'status': 'ok',
                           'code': 400,
                           'error': 'action "%s" invalid' % action
                       }, 400
            return getattr(self, action)(id)

        # list items
        elif id is None:
            # filter
            cond = {}
            if args.get('filter') is not None:
                cond = args.filter
                # cond = json.loads(args.filter)

            # page number
            page = 1
            if args.get('page_num') is not None:
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

            # TODO: getting status for node

            return jsonify({
                'status': 'ok',
                'total_count': total_count,
                'page_num': page,
                'page_size': page_size,
                'items': items
            })

        # get item by id
        else:
            return jsonify(db_manager.get(col_name=self.col_name, id=id))

    def put(self):
        args = self.parser.parse_args()
        item = {}
        for k in args.keys():
            if k not in DEFAULT_ARGS:
                item[k] = args.get(k)
        item = db_manager.save(col_name=self.col_name, item=item)
        return item

    def update(self, id=None):
        args = self.parser.parse_args()
        item = db_manager.get(col_name=self.col_name, id=id)
        if item is None:
            return {
                       'status': 'ok',
                       'code': 401,
                       'error': 'item not exists'
                   }, 401
        values = {}
        for k in args.keys():
            if k not in DEFAULT_ARGS:
                values[k] = args.get(k)
        item = db_manager.update_one(col_name=self.col_name, id=id, values=values)

        # execute after_update hook
        self.after_update(id)

        return item

    def post(self, id=None, action=None):
        if action is None:
            return self.update(id)

        if not hasattr(self, action):
            return {
                       'status': 'ok',
                       'code': 400,
                       'error': 'action "%s" invalid' % action
                   }, 400

        return getattr(self, action)(id)

    def delete(self, id=None):
        db_manager.remove_one(col_name=self.col_name, id=id)

    def after_update(self, id=None):
        pass
