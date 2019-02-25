# -*- coding: utf-8 -*-
import os

import scrapy

from ..items import TaobaoItem


class TaobaoSpiderSpider(scrapy.Spider):
    name = 'taobao_spider'
    allowed_domains = ['taobao.com']
    start_urls = ['http://taobao.com/']

    def parse(self, response):
        yield TaobaoItem()
