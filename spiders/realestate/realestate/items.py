# -*- coding: utf-8 -*-

# Define here the models for your scraped items
#
# See documentation in:
# https://doc.scrapy.org/en/latest/topics/items.html

import scrapy


class RealEstateItem(scrapy.Item):
    # _id
    _id = scrapy.Field()

    # task_id
    task_id = scrapy.Field()

    # 房产名
    name = scrapy.Field()

    # url
    url = scrapy.Field()

    # 类别
    type = scrapy.Field()

    # 价格（万）
    price = scrapy.Field()

    # 大小
    size = scrapy.Field()

    # 小区
    region = scrapy.Field()

    # 城市
    city = scrapy.Field()
