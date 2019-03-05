from utils import node
from .celery import celery_app


@celery_app.task
def update_node_status():
    node.update_nodes_status(refresh=True)
