import os
import subprocess
import sys
from multiprocessing import Process

import click
from flask import Flask
from flask_cors import CORS
from flask_restful import Api

from routes.schedules import ScheduleApi
from tasks.scheduler import scheduler

file_dir = os.path.dirname(os.path.realpath(__file__))
root_path = os.path.abspath(os.path.join(file_dir, '.'))
sys.path.append(root_path)

from config import FLASK_HOST, FLASK_PORT, PROJECT_LOGS_FOLDER, BROKER_URL
from constants.manage import ActionType
from routes.deploys import DeployApi
from routes.files import FileApi
from routes.nodes import NodeApi
from routes.spiders import SpiderApi, SpiderImportApi, SpiderManageApi
from routes.stats import StatsApi
from routes.tasks import TaskApi
from tasks.celery import celery_app

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
api.add_resource(SpiderImportApi,
                 '/api/spiders/import/<string:platform>')
api.add_resource(SpiderManageApi,
                 '/api/spiders/manage/<string:action>')
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
api.add_resource(ScheduleApi,
                 '/api/schedules',
                 '/api/schedules/<string:id>')

if __name__ == '__main__':
    # create folder if it does not exist
    if not os.path.exists(PROJECT_LOGS_FOLDER):
        os.makedirs(PROJECT_LOGS_FOLDER)

    # run app instance
    app.run(host=FLASK_HOST, port=FLASK_PORT, threaded=True)
