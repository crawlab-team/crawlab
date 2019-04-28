import os
from collections import defaultdict
from datetime import datetime, timedelta

from flask_restful import reqparse, Resource

from constants.task import TaskStatus
from db.manager import db_manager
from routes.base import BaseApi
from utils import jsonify


class StatsApi(BaseApi):
    arguments = [
        ['spider_id', str],
    ]

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

    def get_spider_stats(self):
        args = self.parser.parse_args()
        spider_id = args.get('spider_id')
        spider = db_manager.get('spiders', id=spider_id)
        tasks = db_manager.list(
            col_name='tasks',
            cond={
                'spider_id': spider['_id'],
                'create_ts': {
                    '$gte': datetime.now() - timedelta(30)
                }
            },
            limit=9999999
        )

        # task count
        task_count = len(tasks)

        # calculate task count stats
        task_count_by_status = defaultdict(int)
        task_count_by_node = defaultdict(int)
        total_seconds = 0
        for task in tasks:
            task_count_by_status[task['status']] += 1
            task_count_by_node[task.get('node_id')] += 1
            if task['status'] == TaskStatus.SUCCESS and task.get('finish_ts'):
                duration = (task['finish_ts'] - task['create_ts']).total_seconds()
                total_seconds += duration

        # task count by node
        task_count_by_node_ = []
        for status, value in task_count_by_node.items():
            task_count_by_node_.append({
                'name': status,
                'value': value
            })

        # task count by status
        task_count_by_status_ = []
        for status, value in task_count_by_status.items():
            task_count_by_status_.append({
                'name': status,
                'value': value
            })

        # success rate
        success_rate = task_count_by_status[TaskStatus.SUCCESS] / task_count

        # average duration
        avg_duration = total_seconds / task_count

        # calculate task count by date
        cur = db_manager.aggregate('tasks', [
            {
                '$match': {
                    'spider_id': spider['_id']
                }
            },
            {
                '$project': {
                    'date': {
                        '$dateToString': {
                            'format': '%Y-%m-%d',
                            'date': '$create_ts'
                        }
                    },
                    'duration': {
                        '$subtract': [
                            '$finish_ts',
                            '$create_ts'
                        ]
                    }
                }
            },
            {
                '$group': {
                    '_id': '$date',
                    'count': {
                        '$sum': 1
                    },
                    'duration': {
                        '$avg': '$duration'
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
            date_cache[item['_id']] = {
                'duration': item['duration'] / 1000,
                'count': item['count']
            }
        start_date = datetime.now() - timedelta(31)
        end_date = datetime.now() - timedelta(1)
        date = start_date
        daily_tasks = []
        while date < end_date:
            date = date + timedelta(1)
            date_str = date.strftime('%Y-%m-%d')
            d = date_cache.get(date_str)
            row = {
                'date': date_str,
            }
            if d is None:
                row['count'] = 0
                row['duration'] = 0
            else:
                row['count'] = d['count']
                row['duration'] = d['duration']
            daily_tasks.append(row)

        # calculate total results
        result_count = 0
        col_name = spider.get('col')
        if col_name is not None:
            for task in tasks:
                result_count += db_manager.count(col_name, {'task_id': task['_id']})

        # top tasks
        # top_10_tasks = db_manager.list('tasks', {'spider_id': spider['_id']})

        return {
            'status': 'ok',
            'overview': {
                'task_count': task_count,
                'result_count': result_count,
                'success_rate': success_rate,
                'avg_duration': avg_duration
            },
            'task_count_by_status': task_count_by_status_,
            'task_count_by_node': task_count_by_node_,
            'daily_stats': daily_tasks,
        }
