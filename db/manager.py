from bson import ObjectId
from mongoengine import connect
from pymongo import MongoClient, DESCENDING
from config.db import MONGO_HOST, MONGO_PORT, MONGO_DB
from utils import is_object_id

connect(db=MONGO_DB, host=MONGO_HOST, port=MONGO_PORT)


class DbManager(object):
    def __init__(self):
        self.mongo = MongoClient(host=MONGO_HOST, port=MONGO_PORT)
        self.db = self.mongo[MONGO_DB]

    def save(self, col_name: str, item, **kwargs):
        col = self.db[col_name]
        col.save(item, **kwargs)

    def remove(self, col_name: str, cond: dict, **kwargs):
        col = self.db[col_name]
        col.remove(cond, **kwargs)

    def update(self, col_name: str, cond: dict, values: dict, **kwargs):
        col = self.db[col_name]
        col.update(cond, {'$set': values}, **kwargs)

    def update_one(self, col_name: str, id: str, values: dict, **kwargs):
        col = self.db[col_name]
        col.find_one_and_update({'_id': ObjectId(id)}, {'$set': values})

    def remove_one(self, col_name: str, id: str, **kwargs):
        col = self.db[col_name]
        col.remove({'_id': ObjectId(id)})

    def list(self, col_name: str, cond: dict, skip: int = 0, limit: int = 10, **kwargs):
        col = self.db[col_name]
        data = []
        for item in col.find(cond).skip(skip).limit(limit):
            data.append(item)
        return data

    def get(self, col_name: str, id: str):
        if is_object_id(id):
            _id = ObjectId(id)
        else:
            _id = id
        col = self.db[col_name]
        print(_id)
        return col.find_one({'_id': _id})

    def count(self, col_name: str, cond):
        col = self.db[col_name]
        return col.count(cond)

    def get_latest_version(self, spider_id):
        col = self.db['deploys']
        for item in col.find({'spider_id': ObjectId(spider_id)}).sort('version', DESCENDING):
            return item.get('version')
        return None


db_manager = DbManager()
