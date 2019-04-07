import json

import requests

from config import FLOWER_API_ENDPOINT
from constants.node import NodeStatus
from db.manager import db_manager


def check_nodes_status():
    """
    Update node status from Flower.
    """
    res = requests.get('%s/workers?status=1' % FLOWER_API_ENDPOINT)
    return json.loads(res.content.decode('utf-8'))


def update_nodes_status(refresh=False):
    """
    Update all nodes status
    :param refresh:
    """
    online_node_ids = []
    url = '%s/workers?status=1' % FLOWER_API_ENDPOINT
    if refresh:
        url += '&refresh=1'
    res = requests.get(url)
    for k, v in json.loads(res.content.decode('utf-8')).items():
        node_name = k
        node_status = NodeStatus.ONLINE if v else NodeStatus.OFFLINE
        # node_celery = v
        node = db_manager.get('nodes', id=node_name)

        # new node
        if node is None:
            node = {'_id': node_name, 'name': node_name, 'status': node_status}
            db_manager.save('nodes', node)

        else:
            node['status'] = node_status
            db_manager.save('nodes', node)

        if node_status:
            online_node_ids.append(node_name)
    return online_node_ids
