import os

from flask_restful import reqparse, Resource

from utils import jsonify
from utils.file import get_file_content


class FileApi(Resource):
    parser = reqparse.RequestParser()
    arguments = []

    def __init__(self):
        super(FileApi).__init__()
        self.parser.add_argument('path', type=str)

    def get(self, action=None):
        """
        GET method of FileAPI.
        :param action: action
        """
        args = self.parser.parse_args()
        path = args.get('path')

        if action is not None:
            if action == 'getDefaultPath':
                return jsonify({
                    'defaultPath': os.path.abspath(os.path.join(os.path.curdir, 'spiders'))
                })

            elif action == 'get_file':
                file_data = get_file_content(path)
                file_data['status'] = 'ok'
                return jsonify(file_data)

            else:
                return {}

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
