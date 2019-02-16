from app import api
from routes.base import BaseApi


class NodeApi(BaseApi):
    col_name = 'nodes'

    arguments = (
        ('ip', str),
        ('port', int),
        ('name', str),
        ('description', str),
    )


api.add_resource(NodeApi,
                 '/api/nodes',
                 '/api/nodes/<string:id>',
                 '/api/nodes/<string:id>/<string:action>')
