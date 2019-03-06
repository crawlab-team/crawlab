# -*- coding: utf-8 -*-
from time import sleep

import scrapy


class BaiduSpiderSpider(scrapy.Spider):
    name = 'baidu_spider'
    allowed_domains = ['baidu.com']
    start_urls = ['http://baidu.com/s?wd=百度']

    def parse(self, response):
        sleep(30)
