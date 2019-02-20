# -*- coding: utf-8 -*-
import scrapy


class TaobaoSpiderSpider(scrapy.Spider):
    name = 'taobao_spider'
    allowed_domains = ['taobao.com']
    start_urls = ['http://taobao.com/']

    def parse(self, response):
        pass
