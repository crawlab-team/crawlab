# -*- coding: utf-8 -*-
import scrapy
from config_spider.items import Item


class ConfigSpider(scrapy.Spider):
    name = 'config_spider'

    def start_requests(self):
        return scrapy.Request(url='###START_URL###', callback='###START_STAGE###')

###PARSERS###
