import os
from datetime import datetime, timedelta

from flask_restful import reqparse, Resource

from db.manager import db_manager
from utils import jsonify


class StatsApi(Resource):
    def get(self, action: str = None) -> (dict, tuple):
        """
        GET method of StatsApi.
        :param action: action
        """
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
        """
        Get stats for home page
        """
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
        date_cache = {}
        for item in cur:
            date_cache[item['_id']] = item['count']
        start_date = datetime.now() - timedelta(31)
        end_date = datetime.now() - timedelta(1)
        date = start_date
        daily_tasks = []
        while date < end_date:
            date = date + timedelta(1)
            date_str = date.strftime('%Y-%m-%d')
            daily_tasks.append({
                'date': date_str,
                'count': date_cache.get(date_str) or 0,
            })

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
