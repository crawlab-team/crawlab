import json
import re

from bson import json_util


def is_object_id(id):
    return re.search('^[a-zA-Z0-9]{24}$', id) is not None


def jsonify(obj: dict):
    dump_str = json_util.dumps(obj)
    converted_obj = json.loads(dump_str)
    return converted_obj
