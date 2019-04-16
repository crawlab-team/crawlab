# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: http://doc.scrapy.org/en/latest/topics/item-pipeline.html
import os

from pymongo import MongoClient

MONGO_HOST = os.environ['MONGO_HOST']
MONGO_PORT = int(os.environ['MONGO_PORT'])
MONGO_DB = os.environ['MONGO_DB']


class MongoPipeline(object):
    mongo = MongoClient(host=MONGO_HOST, port=MONGO_PORT)
    db = mongo[MONGO_DB]
    col_name = os.environ.get('CRAWLAB_COLLECTION')
    col = db[col_name]

    def process_item(self, item, spider):
        item['task_id'] = os.environ.get('CRAWLAB_TASK_ID')
        self.col.save(item)
        return item
