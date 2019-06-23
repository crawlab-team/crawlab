# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://doc.scrapy.org/en/latest/topics/item-pipeline.html
import os

from pymongo import MongoClient


class XueqiuPipeline(object):
    mongo = MongoClient(
        host='localhost',
        port=27017
    )
    db = mongo['crawlab_test']
    col = db.get_collection(os.environ.get('CRAWLAB_COLLECTION') or 'results_xueqiu')

    def process_item(self, item, spider):
        item['task_id'] = os.environ.get('CRAWLAB_TASK_ID')
        item['_id'] = item['id']
        if self.col.find_one({'_id': item['_id']}) is None:
            self.col.save(item)
            return item
