from celery import Celery
from flask import Flask
from flask_cors import CORS
from flask_restful import Api, Resource

from routes.deploys import DeployApi
from routes.files import FileApi
from routes.nodes import NodeApi
from routes.spiders import SpiderApi
from routes.tasks import TaskApi

# flask app instance
app = Flask(__name__)
app.config.from_object('config.flask')
# init flask api instance
api = Api(app)

# cors support
CORS(app, supports_credentials=True)

# reference api routes

api.add_resource(NodeApi,
                 '/api/nodes',
                 '/api/nodes/<string:id>',
                 '/api/nodes/<string:id>/<string:action>')
api.add_resource(SpiderApi,
                 '/api/spiders',
                 '/api/spiders/<string:id>',
                 '/api/spiders/<string:id>/<string:action>')
api.add_resource(DeployApi,
                 '/api/deploys',
                 '/api/deploys/<string:id>',
                 '/api/deploys/<string:id>/<string:action>')
api.add_resource(TaskApi,
                 '/api/tasks',
                 '/api/tasks/<string:id>'
                 )
api.add_resource(FileApi,
                 '/api/files',
                 '/api/files/<string:action>')

# start flask app
if __name__ == '__main__':
    app.run()
