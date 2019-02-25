import os

from flask_restful import reqparse, Resource

from db.manager import db_manager
from utils import jsonify


class StatsApi(Resource):
    def get(self, action=None):
        # action
        if action is not None:
            if not hasattr(self, action):
                return {
                           'status': 'ok',
                           'code': 400,
                           'error': 'action "%s" invalid' % action
                       }, 400
            return getattr(self, action)()

        else:
            return {}

    def get_home_stats(self):
        # overview stats
        task_count = db_manager.count('tasks', {})
        spider_count = db_manager.count('spiders', {})
        node_count = db_manager.count('nodes', {})
        deploy_count = db_manager.count('deploys', {})

        # daily stats
        cur = db_manager.aggregate('tasks', [
            {
                '$project': {
                    'date': {
                        '$dateToString': {
                            'format': '%Y-%m-%d',
                            'date': '$create_ts'
                        }
                    }
                }
            },
            {
                '$group': {
                    '_id': '$date',
                    'count': {
                        '$sum': 1
                    }
                }
            },
            {
                '$sort': {
                    '_id': 1
                }
            }
        ])
        daily_tasks = []
        for item in cur:
            daily_tasks.append(item)

        return {
            'status': 'ok',
            'overview_stats': {
                'task_count': task_count,
                'spider_count': spider_count,
                'node_count': node_count,
                'deploy_count': deploy_count,
            },
            'daily_tasks': daily_tasks
        }
