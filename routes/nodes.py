import json

import requests

from app import api
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

    def get(self, id=None):
        if id is not None:
            return db_manager.get('nodes', id=id)

        else:
            res = requests.get('%s/workers' % FLOWER_API_ENDPOINT)
            for k, v in json.loads(res.content.decode('utf-8')).items():
                node_name = k
                node = v
                node['_id'] = node_name
                node['name'] = node_name
                node['status'] = NodeType.ONLINE
                db_manager.save('nodes', node)

        items = db_manager.list('nodes', {})

        return jsonify({
            'status': 'ok',
            'items': items
        })


api.add_resource(NodeApi,
                 '/api/nodes',
                 '/api/nodes/<string:id>',
                 '/api/nodes/<string:id>/<string:action>')
