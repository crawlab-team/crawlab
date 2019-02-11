from celery import Celery
from flask import Flask
from flask_restful import Api

# TODO: 用配置文件启动 http://www.pythondoc.com/flask/config.html
app = Flask(__name__)
app.config['DEBUG'] = True

# init flask api instance
api = Api(app)

# reference api routes
import routes.tasks
import routes.spiders
import routes.test

# start flask app
if __name__ == '__main__':
    app.run()
