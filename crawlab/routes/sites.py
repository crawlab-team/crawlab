import json

from bson import ObjectId
from pymongo import ASCENDING

from db.manager import db_manager
from routes.base import BaseApi
from utils import jsonify


class SiteApi(BaseApi):
    col_name = 'sites'

    arguments = (
        ('keyword', str),
        ('category', str),
    )

    def get(self, id: str = None, action: str = None):
        # action by id
        if action is not None:
            if not hasattr(self, action):
                return {
                           'status': 'ok',
                           'code': 400,
                           'error': 'action "%s" invalid' % action
                       }, 400
            return getattr(self, action)(id)

        elif id is not None:
            site = db_manager.get(col_name=self.col_name, id=id)
            return jsonify(site)

        # list tasks
        args = self.parser.parse_args()
        page_size = args.get('page_size') or 10
        page_num = args.get('page_num') or 1
        filter_str = args.get('filter')
        keyword = args.get('keyword')
        filter_ = {}
        if filter_str is not None:
            filter_ = json.loads(filter_str)
        if keyword is not None:
            filter_['$or'] = [
                {'description': {'$regex': keyword}},
                {'name': {'$regex': keyword}}
            ]

        items = db_manager.list(
            col_name=self.col_name,
            cond=filter_,
            limit=page_size,
            skip=page_size * (page_num - 1),
            sort_key='rank',
            sort_direction=ASCENDING
        )

        return {
            'status': 'ok',
            'total_count': db_manager.count(self.col_name, filter_),
            'page_num': page_num,
            'page_size': page_size,
            'items': jsonify(items)
        }
