import subprocess

from app import api
from config.celery import BROKER_URL
from routes.base import BaseApi


class NodeApi(BaseApi):
    col_name = 'nodes'

    arguments = (
        ('ip', str),
        ('port', int),
        ('name', str),
        ('description', str),
    )

    def _get(self, id=None):
        if id is not None:
            return {
            }

        else:
            p = subprocess.Popen(['celery', 'inspect', 'stats', '-b', BROKER_URL])
            stdout, stderr = p.communicate()
            return {
                'stdout': stdout,
                'stderr': stderr,
            }


api.add_resource(NodeApi,
                 '/api/nodes',
                 '/api/nodes/<string:id>',
                 '/api/nodes/<string:id>/<string:action>')
