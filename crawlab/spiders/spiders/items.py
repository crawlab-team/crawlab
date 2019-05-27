# -*- coding: utf-8 -*-

# Define here the models for your scraped items
#
# See documentation in:
# https://doc.scrapy.org/en/latest/topics/items.html

import scrapy

from spiders.db import spider


class SpidersItem(scrapy.Item):
    if spider['crawl_type'] == 'list':
        fields = {f['name']: scrapy.Field() for f in spider['fields']}
    elif spider['crawl_type'] == 'detail':
        fields = {f['name']: scrapy.Field() for f in spider['detail_fields']}
    elif spider['crawl_type'] == 'list-detail':
        fields = {f['name']: scrapy.Field() for f in (spider['fields'] + spider['detail_fields'])}
    else:
        fields = {}

    # basic fields
    fields['_id'] = scrapy.Field()
    fields['task_id'] = scrapy.Field()
