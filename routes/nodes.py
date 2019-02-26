import json

import requests
from bson import ObjectId

from config.celery import FLOWER_API_ENDPOINT
from constants.node import NodeType
from db.manager import db_manager
from routes.base import BaseApi
from utils import jsonify


class NodeApi(BaseApi):
    col_name = 'nodes'

    arguments = (
        # ('ip', str),
        # ('port', int),
        ('name', str),
        ('description', str),
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

        # get one node
        elif id is not None:
            return db_manager.get('nodes', id=id)

        # TODO: use query "?status=1" to get status of nodes

        # get a list of items
        res = requests.get('%s/workers' % FLOWER_API_ENDPOINT)
        online_node_ids = []
        for k, v in json.loads(res.content.decode('utf-8')).items():
            node_name = k
            node_celery = v
            node = db_manager.get('nodes', id=node_name)

            # new node
            if node is None:
                node = {}
                for _k, _v in node_celery.items():
                    node[_k] = _v
                node['_id'] = node_name
                node['name'] = node_name
                node['status'] = NodeType.ONLINE
                db_manager.save('nodes', node)

            # existing node
            else:
                for _k, _v in v.items():
                    node[_k] = _v
                node['name'] = node_name
                node['status'] = NodeType.ONLINE
                db_manager.save('nodes', node)

            online_node_ids.append(node_name)

        # iterate db nodes to update status
        nodes = []
        items = db_manager.list('nodes', {})
        for item in items:
            if item['_id'] in online_node_ids:
                item['status'] = NodeType.ONLINE
            else:
                item['status'] = NodeType.OFFLINE
            db_manager.update_one('nodes', item['_id'], {
                'status': item['status']
            })
            nodes.append(item)

        return jsonify({
            'status': 'ok',
            'items': nodes
        })

    def get_spiders(self, id=None):
        items = db_manager.list('spiders')

    def get_deploys(self, id):
        items = db_manager.list('deploys', {'node_id': id}, limit=10, sort_key='finish_ts')
        deploys = []
        for item in items:
            spider_id = item['spider_id']
            spider = db_manager.get('spiders', id=str(spider_id))
            item['spider_name'] = spider['name']
            deploys.append(item)
        return jsonify({
            'status': 'ok',
            'items': deploys
        })

    def get_tasks(self, id):
        items = db_manager.list('tasks', {'node_id': id}, limit=10, sort_key='create_ts')
        for item in items:
            spider_id = item['spider_id']
            spider = db_manager.get('spiders', id=str(spider_id))
            item['spider_name'] = spider['name']
            task = db_manager.get('tasks_celery', id=item['_id'])
            item['status'] = task['status']
        return jsonify({
            'status': 'ok',
            'items': items
        })
