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
    """
    Base class for API. All API classes should inherit this class.
    """
    col_name = 'tmp'
    parser = reqparse.RequestParser()
    arguments = []

    def __init__(self):
        super(BaseApi).__init__()
        self.parser.add_argument('page_num', type=int)
        self.parser.add_argument('page_size', type=int)
        self.parser.add_argument('filter', type=str)

        for arg, type in self.arguments:
            self.parser.add_argument(arg, type=type)

    def get(self, id: str = None, action: str = None) -> (dict, tuple):
        """
        GET method for retrieving item information.
        If id is specified and action is not, return the object of the given id;
        If id and action are both specified, execute the given action results of the given id;
        If neither id nor action is specified, return the list of items given the page_size, page_num and filter
        :param id:
        :param action:
        :return:
        """
        # import pdb
        # pdb.set_trace()
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

            return {
                'status': 'ok',
                'total_count': total_count,
                'page_num': page,
                'page_size': page_size,
                'items': jsonify(items)
            }

        # get item by id
        else:
            return jsonify(db_manager.get(col_name=self.col_name, id=id))

    def put(self) -> (dict, tuple):
        """
        PUT method for creating a new item.
        :return:
        """
        args = self.parser.parse_args()
        item = {}
        for k in args.keys():
            if k not in DEFAULT_ARGS:
                item[k] = args.get(k)
        id = db_manager.save(col_name=self.col_name, item=item)

        # execute after_update hook
        self.after_update(id)

        return jsonify(id)

    def update(self, id: str = None) -> (dict, tuple):
        """
        Helper function for update action given the id.
        :param id:
        :return:
        """
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
                if args.get(k) is not None:
                    values[k] = args.get(k)
        item = db_manager.update_one(col_name=self.col_name, id=id, values=values)

        # execute after_update hook
        self.after_update(id)

        return jsonify(item)

    def post(self, id: str = None, action: str = None):
        """
        POST method of the given id for performing an action.
        :param id:
        :param action:
        :return:
        """
        # perform update action if action is not specified
        if action is None:
            return self.update(id)

        # if action is not defined in the attributes, return 400 error
        if not hasattr(self, action):
            return {
                       'status': 'ok',
                       'code': 400,
                       'error': 'action "%s" invalid' % action
                   }, 400

        # perform specified action of given id
        return getattr(self, action)(id)

    def delete(self, id: str = None) -> (dict, tuple):
        """
        DELETE method of given id for deleting an item.
        :param id:
        :return:
        """
        # perform delete action
        db_manager.remove_one(col_name=self.col_name, id=id)

        # execute after_update hook
        self.after_update(id)

        return {
            'status': 'ok',
            'message': 'deleted successfully',
        }

    def after_update(self, id: str = None):
        """
        This is the after update hook once the update method is performed.
        To be overridden.
        :param id:
        :return:
        """
        pass
