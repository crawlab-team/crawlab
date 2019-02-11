import re


def is_object_id(id):
    return re.search('^[a-zA-Z0-9]{24}$', id) is not None
