import os

from pymongo import MongoClient

MONGO_HOST = os.environ.get('MONGO_HOST')
MONGO_PORT = int(os.environ.get('MONGO_PORT'))
MONGO_DB = os.environ.get('MONGO_DB')
mongo = MongoClient(host=MONGO_HOST,
                    port=MONGO_PORT)
db = mongo[MONGO_DB]
task_id = os.environ.get('CRAWLAB_TASK_ID')
col_name = os.environ.get('CRAWLAB_COLLECTION')
task = db['tasks'].find_one({'_id': task_id})
spider = db['spiders'].find_one({'_id': task['spider_id']})
