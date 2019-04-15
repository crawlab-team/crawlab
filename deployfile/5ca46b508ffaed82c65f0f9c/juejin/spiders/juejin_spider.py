# -*- coding: utf-8 -*-
import scrapy
from juejin.items import JuejinItem


class JuejinSpiderSpider(scrapy.Spider):
    name = 'juejin_spider'
    allowed_domains = ['juejin.com']
    start_urls = ['https://juejin.im/search?query=celery']

    def parse(self, response):
        for item in response.css('ul.main-list > li.item'):
            yield JuejinItem(
                title=item.css('.title span').extract_first(),
                link=item.css('a::attr("href")').extract_first(),
                like=item.css('.like .count::text').extract_first(),
            )
