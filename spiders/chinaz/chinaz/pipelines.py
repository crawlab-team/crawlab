# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://doc.scrapy.org/en/latest/topics/item-pipeline.html

import os

from pymongo import MongoClient

MONGO_HOST = os.environ.get('MONGO_HOST') or 'localhost'
MONGO_PORT = int(os.environ.get('MONGO_PORT') or '27017')
MONGO_DB = os.environ.get('MONGO_DB') or 'crawlab_test'


class MongoPipeline(object):
    mongo = MongoClient(host=MONGO_HOST, port=MONGO_PORT)
    db = mongo[MONGO_DB]
    col_name = os.environ.get('CRAWLAB_COLLECTION') or 'sites'
    col = db[col_name]

    def process_item(self, item, spider):
        item['task_id'] = os.environ.get('CRAWLAB_TASK_ID')
        item['_id'] = item['domain']
        if self.col.find_one({'_id': item['_id']}) is None:
            self.col.save(item)
        return item
