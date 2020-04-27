# -*- coding: utf-8 -*-

# Define here the models for your scraped items
#
# See documentation in:
# https://docs.scrapy.org/en/latest/topics/items.html

import scrapy


class Item(scrapy.Item):
    _id = scrapy.Field()
    task_id = scrapy.Field()
    ts = scrapy.Field()
    title = scrapy.Field()
    url = scrapy.Field()
    abstract = scrapy.Field()
    votes = scrapy.Field()

