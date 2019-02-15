from app import api
from routes.base import BaseApi


class TaskApi(BaseApi):
    col_name = 'tasks_celery'

    arguments = (
        ('deploy_id', str),
        ('file_path', str)
    )


# add api to resources
api.add_resource(TaskApi,
                 '/api/tasks',
                 '/api/tasks/<string:id>'
                 )
