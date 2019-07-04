# -*- coding: utf-8 -*-

# Define here the models for your scraped items
#
# See documentation in:
# https://doc.scrapy.org/en/latest/topics/items.html

import scrapy


class NewsItem(scrapy.Item):
    # define the fields for your item here like:
    _id = scrapy.Field()
    title = scrapy.Field()
    ts_str = scrapy.Field()
    ts = scrapy.Field()
    url = scrapy.Field()
    text = scrapy.Field()
    task_id = scrapy.Field()
    source = scrapy.Field()
    stocks = scrapy.Field()
