from bson import ObjectId
from mongoengine import connect
from pymongo import MongoClient, DESCENDING
from config import MONGO_HOST, MONGO_PORT, MONGO_DB
from utils import is_object_id, jsonify

connect(db=MONGO_DB, host=MONGO_HOST, port=MONGO_PORT)


class DbManager(object):
    def __init__(self):
        self.mongo = MongoClient(host=MONGO_HOST, port=MONGO_PORT)
        self.db = self.mongo[MONGO_DB]

    def save(self, col_name: str, item, **kwargs):
        col = self.db[col_name]

        # in case some fields cannot be saved in MongoDB
        if item.get('stats') is not None:
            item.pop('stats')

        col.save(item, **kwargs)

    def remove(self, col_name: str, cond: dict, **kwargs):
        col = self.db[col_name]
        col.remove(cond, **kwargs)

    def update(self, col_name: str, cond: dict, values: dict, **kwargs):
        col = self.db[col_name]
        col.update(cond, {'$set': values}, **kwargs)

    def update_one(self, col_name: str, id: str, values: dict, **kwargs):
        col = self.db[col_name]
        _id = id
        if is_object_id(id):
            _id = ObjectId(id)
        # print('UPDATE: _id = "%s", values = %s' % (str(_id), jsonify(values)))
        col.find_one_and_update({'_id': _id}, {'$set': values})

    def remove_one(self, col_name: str, id: str, **kwargs):
        col = self.db[col_name]
        _id = id
        if is_object_id(id):
            _id = ObjectId(id)
        col.remove({'_id': _id})

    def list(self, col_name: str, cond: dict, sort_key=None, sort_direction=DESCENDING, skip: int = 0, limit: int = 100,
             **kwargs):
        if sort_key is None:
            sort_key = '_i'
        col = self.db[col_name]
        data = []
        for item in col.find(cond).sort(sort_key, sort_direction).skip(skip).limit(limit):
            data.append(item)
        return data

    def _get(self, col_name: str, cond: dict):
        col = self.db[col_name]
        return col.find_one(cond)

    def get(self, col_name: str, id):
        if type(id) == ObjectId:
            _id = id
        elif is_object_id(id):
            _id = ObjectId(id)
        else:
            _id = id
        return self._get(col_name=col_name, cond={'_id': _id})

    def get_one_by_key(self, col_name: str, key, value):
        return self._get(col_name=col_name, cond={key: value})

    def count(self, col_name: str, cond):
        col = self.db[col_name]
        return col.count(cond)

    def get_latest_version(self, spider_id, node_id):
        col = self.db['deploys']
        for item in col.find({'spider_id': ObjectId(spider_id), 'node_id': node_id}) \
                .sort('version', DESCENDING):
            return item.get('version')
        return None

    def get_last_deploy(self, spider_id):
        col = self.db['deploys']
        for item in col.find({'spider_id': ObjectId(spider_id)}) \
                .sort('finish_ts', DESCENDING):
            return item
        return None

    def aggregate(self, col_name: str, pipelines, **kwargs):
        col = self.db[col_name]
        return col.aggregate(pipelines, **kwargs)


db_manager = DbManager()
