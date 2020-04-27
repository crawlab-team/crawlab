# -*- coding: utf-8 -*-

# Define here the models for your scraped items
#
# See documentation in:
# https://doc.scrapy.org/en/latest/topics/items.html

import scrapy


class ChinazItem(scrapy.Item):
    # define the fields for your item here like:
    _id = scrapy.Field()
    task_id = scrapy.Field()
    name = scrapy.Field()
    domain = scrapy.Field()
    description = scrapy.Field()
    rank = scrapy.Field()
    main_category = scrapy.Field()
    category = scrapy.Field()
    location = scrapy.Field()
