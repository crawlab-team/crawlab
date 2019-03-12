import json

import requests
from celery.worker.control import revoke

from constants.task import TaskStatus
from db.manager import db_manager
from routes.base import BaseApi
from utils import jsonify
from utils.spider import get_spider_col_fields


class TaskApi(BaseApi):
    col_name = 'tasks'

    arguments = (
        ('deploy_id', str),
        ('file_path', str)
    )

    def get(self, id=None, action=None):
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
            task = db_manager.get('tasks', id=id)
            _task = db_manager.get('tasks_celery', id=task['_id'])
            _spider = db_manager.get('spiders', id=str(task['spider_id']))
            if _task:
                if not task.get('status'):
                    task['status'] = _task['status']
            # else:
            #     task['status'] = TaskStatus.UNAVAILABLE
            task['result'] = _task['result']
            task['spider_name'] = _spider['name']
            try:
                with open(task['log_file_path']) as f:
                    task['log'] = f.read()
            except Exception as err:
                task['log'] = ''
            return jsonify(task)

        tasks = db_manager.list('tasks', {}, limit=1000, sort_key='finish_ts')
        items = []
        for task in tasks:
            _task = db_manager.get('tasks_celery', id=task['_id'])
            _spider = db_manager.get('spiders', id=str(task['spider_id']))
            if _task:
                task['status'] = _task['status']
            else:
                task['status'] = TaskStatus.UNAVAILABLE
            task['spider_name'] = _spider['name']
            items.append(task)
        return jsonify({
            'status': 'ok',
            'items': items
        })

    def on_get_log(self, id):
        try:
            task = db_manager.get('tasks', id=id)
            with open(task['log_file_path']) as f:
                log = f.read()
                return {
                    'status': 'ok',
                    'log': log
                }
        except Exception as err:
            return {
                       'code': 500,
                       'status': 'ok',
                       'error': str(err)
                   }, 500

    def get_log(self, id):
        task = db_manager.get('tasks', id=id)
        node = db_manager.get('nodes', id=task['node_id'])
        r = requests.get('http://%s:%s/api/tasks/%s/on_get_log' % (
            node['ip'],
            node['port'],
            id
        ))
        if r.status_code == 200:
            data = json.loads(r.content.decode('utf-8'))
            return {
                'status': 'ok',
                'log': data.get('log')
            }
        else:
            data = json.loads(r.content)
            return {
                       'code': 500,
                       'status': 'ok',
                       'error': data['error']
                   }, 500

    def get_results(self, id):
        task = db_manager.get('tasks', id=id)
        spider = db_manager.get('spiders', id=task['spider_id'])
        col_name = spider.get('col')
        if not col_name:
            return []
        fields = get_spider_col_fields(col_name)
        items = db_manager.list(col_name, {'task_id': id})
        return jsonify({
            'status': 'ok',
            'fields': fields,
            'items': items
        })

    def stop(self, id):
        revoke(id, terminate=True)
        return {
            'id': id,
            'status': 'ok',
        }
