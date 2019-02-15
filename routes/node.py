from app import api
from routes.base import BaseApi


class NodeApi(BaseApi):
    col_name = 'nodes'

    arguments = (
        ('node_ip', str),
        ('node_name', str),
        ('node_description', str),
    )


api.add_resource(NodeApi,
                 '/api/nodes',
                 '/api/nodes/<string:id>',
                 '/api/nodes/<string:id>/<string:action>')
