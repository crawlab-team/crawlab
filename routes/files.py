import os

from flask_restful import reqparse, Resource

from app import api
from utils import jsonify


class FileApi(Resource):
    parser = reqparse.RequestParser()
    arguments = []

    def __init__(self):
        super(FileApi).__init__()
        self.parser.add_argument('path', type=str)

    def get(self, action=None):
        if action is not None:
            if action == 'getDefaultPath':
                return jsonify({
                    'defaultPath': os.path.abspath(os.path.join(os.path.curdir, 'spiders'))
                })
            else:
                return {}

        args = self.parser.parse_args()
        path = args.get('path')
        folders = []
        files = []
        for _path in os.listdir(path):
            if os.path.isfile(os.path.join(path, _path)):
                files.append(_path)
            elif os.path.isdir(os.path.join(path, _path)):
                folders.append(_path)
        return jsonify({
            'status': 'ok',
            'files': sorted(files),
            'folders': sorted(folders),
        })


api.add_resource(FileApi,
                 '/api/files',
                 '/api/files/<string:action>')
