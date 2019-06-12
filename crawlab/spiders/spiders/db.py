import os

from pymongo import MongoClient

MONGO_HOST = os.environ.get('MONGO_HOST') or 'localhost'
MONGO_PORT = int(os.environ.get('MONGO_PORT')) or 27017
MONGO_USERNAME = os.environ.get('MONGO_USERNAME')
MONGO_PASSWORD = os.environ.get('MONGO_PASSWORD')
MONGO_DB = os.environ.get('MONGO_DB') or 'crawlab_test'
mongo = MongoClient(host=MONGO_HOST,
                    port=MONGO_PORT,
                    username=MONGO_USERNAME,
                    password=MONGO_PASSWORD)
db = mongo[MONGO_DB]
task_id = os.environ.get('CRAWLAB_TASK_ID')
col_name = os.environ.get('CRAWLAB_COLLECTION')
task = db['tasks'].find_one({'_id': task_id})
spider = db['spiders'].find_one({'_id': task['spider_id']})
