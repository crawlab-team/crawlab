import subprocess
import sys
import threading

from celery import Celery
from flask import Flask
from flask_cors import CORS
from flask_restful import Api, Resource

from config import BROKER_URL
from routes.deploys import DeployApi
from routes.files import FileApi
from routes.nodes import NodeApi
from routes.spiders import SpiderApi
from routes.stats import StatsApi
from routes.tasks import TaskApi

# flask app instance
app = Flask(__name__)
app.config.from_object('config')
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
                 '/api/tasks/<string:id>',
                 '/api/tasks/<string:id>/<string:action>'
                 )
api.add_resource(FileApi,
                 '/api/files',
                 '/api/files/<string:action>')
api.add_resource(StatsApi,
                 '/api/stats',
                 '/api/stats/<string:action>')


def run_app():
    app.run()


def run_flower():
    p = subprocess.Popen(['celery', 'flower', '-b', BROKER_URL], stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    for line in iter(p.stdout.readline, 'b'):
        print(line.decode('utf-8'))


if __name__ == '__main__':
    # start flower app
    th_flower = threading.Thread(target=run_flower)
    th_flower.start()

    # start flask app
    # th_app = threading.Thread(target=run_app)
    # th_app.start()
    app.run()

