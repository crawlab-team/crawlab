from routes.base import BaseApi


class DeployApi(BaseApi):
    col_name = 'deploys'

    arguments = (
        ('spider_id', str),
        ('node_id', str),
    )

