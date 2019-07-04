# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://doc.scrapy.org/en/latest/topics/item-pipeline.html
import os

from pymongo import MongoClient


class XueqiuPipeline(object):
    mongo = MongoClient(
        host=os.environ.get('MONGO_HOST') or 'localhost',
        port=int(os.environ.get('MONGO_PORT') or 27017)
    )
    db = mongo[os.environ.get('MONGO_DB') or 'crawlab_test']
    col = db.get_collection(os.environ.get('CRAWLAB_COLLECTION') or 'results_xueqiu')

    # create indexes
    col.create_index('stocks')
    col.create_index('id')
    col.create_index('url')

    def process_item(self, item, spider):
        item['task_id'] = os.environ.get('CRAWLAB_TASK_ID')
        if self.col.find_one({'id': item['id']}) is None:
            self.col.save(item)
            return item
