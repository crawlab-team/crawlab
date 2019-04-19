import os
import subprocess
import sys
from multiprocessing import Process

import click
from celery import Celery
from flask import Flask
from flask_cors import CORS
from flask_restful import Api
# from flask_restplus import Api
from utils.log import other
from constants.node import NodeStatus
from db.manager import db_manager
from routes.schedules import ScheduleApi
from tasks.celery import celery_app
from tasks.scheduler import scheduler

file_dir = os.path.dirname(os.path.realpath(__file__))
root_path = os.path.abspath(os.path.join(file_dir, '.'))
sys.path.append(root_path)

from config import FLASK_HOST, FLASK_PORT, PROJECT_LOGS_FOLDER
from routes.deploys import DeployApi
from routes.files import FileApi
from routes.nodes import NodeApi
from routes.spiders import SpiderApi, SpiderImportApi, SpiderManageApi
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
api.add_resource(SpiderImportApi,
                 '/api/spiders/import/<string:platform>')
api.add_resource(SpiderManageApi,
                 '/api/spiders/manage/<string:action>')
api.add_resource(TaskApi,
                 '/api/tasks',
                 '/api/tasks/<string:id>',
                 '/api/tasks/<string:id>/<string:action>')
api.add_resource(DeployApi,
                 '/api/deploys',
                 '/api/deploys/<string:id>',
                 '/api/deploys/<string:id>/<string:action>')
api.add_resource(FileApi,
                 '/api/files',
                 '/api/files/<string:action>')
api.add_resource(StatsApi,
                 '/api/stats',
                 '/api/stats/<string:action>')
api.add_resource(ScheduleApi,
                 '/api/schedules',
                 '/api/schedules/<string:id>')


def monitor_nodes_status(celery_app):
    def update_nodes_status(event):
        node_id = event.get('hostname')
        db_manager.update_one('nodes', id=node_id, values={
            'status': NodeStatus.ONLINE
        })

    def update_nodes_status_online(event):
        other.info(f"{event}")

    with celery_app.connection() as connection:
        recv = celery_app.events.Receiver(connection, handlers={
            'worker-heartbeat': update_nodes_status,
            # 'worker-online': update_nodes_status_online,
        })
        recv.capture(limit=None, timeout=None, wakeup=True)


# run scheduler as a separate process
scheduler.run()

# monitor node status
p_monitor = Process(target=monitor_nodes_status, args=(celery_app,))
p_monitor.start()

# create folder if it does not exist
if not os.path.exists(PROJECT_LOGS_FOLDER):
    os.makedirs(PROJECT_LOGS_FOLDER)

if __name__ == '__main__':
    # run app instance
    app.run(host=FLASK_HOST, port=FLASK_PORT, threaded=True)
