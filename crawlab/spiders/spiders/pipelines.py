# -*- coding: utf-8 -*-

# Define your item pipelines here
#
# Don't forget to add your pipeline to the ITEM_PIPELINES setting
# See: https://doc.scrapy.org/en/latest/topics/item-pipeline.html
from spiders.db import db, col_name, task_id


class SpidersPipeline(object):
    col = db[col_name]

    def process_item(self, item, spider):
        item['task_id'] = task_id
        self.col.save(item)

        return item
