import json

import requests

from config.celery import FLOWER_API_ENDPOINT
from constants.node import NodeType
from db.manager import db_manager
from routes.base import BaseApi
from utils import jsonify


class NodeApi(BaseApi):
    col_name = 'nodes'

    arguments = (
        ('ip', str),
        ('port', int),
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

        # get a list of items
        else:
            res = requests.get('%s/workers' % FLOWER_API_ENDPOINT)
            for k, v in json.loads(res.content.decode('utf-8')).items():
                node_name = k
                node_celery = v
                node = db_manager.get('nodes', id=node_name)

                # new node
                if node is None:
                    node = {}
                    for _k, _v in v.items():
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

        items = db_manager.list('nodes', {})

        return jsonify({
            'status': 'ok',
            'items': items
        })

    def spider(self, id=None):
        items = db_manager.list('spiders')
