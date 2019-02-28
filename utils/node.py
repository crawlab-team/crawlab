import json

import requests

from config import FLOWER_API_ENDPOINT


def check_nodes_status():
    res = requests.get('%s/workers?status=1' % FLOWER_API_ENDPOINT)
    return json.loads(res.content.decode('utf-8'))
