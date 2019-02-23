from db.manager import db_manager
from routes.base import BaseApi
from utils import jsonify


class TaskApi(BaseApi):
    col_name = 'tasks'

    arguments = (
        ('deploy_id', str),
        ('file_path', str)
    )

    def get(self, id=None):
        if id is not None:
            task = db_manager.get('tasks', id=id)
            _task = db_manager.get('tasks_celery', id=task['_id'])
            _spider = db_manager.get('spiders', id=str(task['spider_id']))
            task['status'] = _task['status']
            task['result'] = _task['result']
            task['spider_name'] = _spider['name']
            with open(task['log_file_path']) as f:
                task['log'] = f.read()
            return jsonify(task)

        tasks = db_manager.list('tasks', {}, limit=1000)
        items = []
        for task in tasks:
            _task = db_manager.get('tasks_celery', id=task['_id'])
            _spider = db_manager.get('spiders', id=str(task['spider_id']))
            task['status'] = _task['status']
            task['spider_name'] = _spider['name']
            items.append(task)
        return jsonify({
            'status': 'ok',
            'items': items
        })
