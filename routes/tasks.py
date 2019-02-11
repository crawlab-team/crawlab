from app import api
from routes.base import BaseApi


class TaskApi(BaseApi):
    col_name = 'tasks_celery'


# add api to resources
api.add_resource(TaskApi,
                 '/api/tasks',
                 '/api/tasks/<string:id>'
                 )
