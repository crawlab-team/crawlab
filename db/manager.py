from pymongo import MongoClient
from config.db import MONGO_HOST, MONGO_PORT, MONGO_DB


class DbManager(object):
    def __init__(self):
        self.mongo = MongoClient(host=MONGO_HOST, port=MONGO_PORT)
        self.db = self.mongo[MONGO_DB]

    # TODO: CRUD
    def save(self, col_name: str, item, **kwargs):
        col = self.db[col_name]
        col.save(item, **kwargs)

    def remove(self, col_name: str, cond: dict, **kwargs):
        col = self.db[col_name]
        col.remove(cond, **kwargs)

    def update(self, col_name: str, cond: dict, values: dict, **kwargs):
        col = self.db[col_name]
        col.update(cond, {'$set': values}, **kwargs)

    def list(self, col_name: str, cond: dict, skip: int, limit: int, **kwargs):
        if kwargs.get('page') is not None:
            try:
                page = int(kwargs.get('page'))
                skip = page * limit
            except Exception as err:
                pass
        # TODO: list logic
        # TODO: pagination
