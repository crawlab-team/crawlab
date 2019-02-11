from celery import Celery
from flask import Flask
from flask_restful import Api

# TODO: 用配置文件启动 http://www.pythondoc.com/flask/config.html
app = Flask(__name__)
app.config['DEBUG'] = True

# init flask api instance
api = Api(app)

# reference api routes
import routes.task
import routes.spider

# start flask app
app.run()
