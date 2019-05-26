# -*- coding: utf-8 -*-

# Define here the models for your scraped items
#
# See documentation in:
# https://doc.scrapy.org/en/latest/topics/items.html

import scrapy

from spiders.db import spider


class SpidersItem(scrapy.Item):
    fields = {f['name']: scrapy.Field() for f in spider['fields']}
    fields['_id'] = scrapy.Field()
    fields['task_id'] = scrapy.Field()
