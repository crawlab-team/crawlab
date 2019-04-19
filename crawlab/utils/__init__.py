import json
import re
from datetime import datetime

from bson import json_util


def is_object_id(id: str) -> bool:
    """
    Determine if the id is a valid ObjectId string
    :param id: ObjectId string
    """
    return re.search('^[a-zA-Z0-9]{24}$', id) is not None


def jsonify(obj: (dict, list)) -> (dict, list):
    """
    Convert dict/list to a valid json object.
    :param obj: object to be converted
    :return: dict/list
    """
    dump_str = json_util.dumps(obj)
    converted_obj = json.loads(dump_str)
    if type(converted_obj) == dict:
        for k, v in converted_obj.items():
            if type(v) == dict:
                if v.get('$oid') is not None:
                    converted_obj[k] = v['$oid']
                elif v.get('$date') is not None:
                    converted_obj[k] = datetime.fromtimestamp(v['$date'] / 1000).strftime('%Y-%m-%d %H:%M:%S')
    elif type(converted_obj) == list:
        for i, v in enumerate(converted_obj):
            converted_obj[i] = jsonify(v)
    return converted_obj
