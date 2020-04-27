# -*- coding: utf-8 -*-
import scrapy
import re
from config_spider.items import Item
from urllib.parse import urljoin, urlparse

def get_real_url(response, url):
    if re.search(r'^https?', url):
        return url
    elif re.search(r'^\/\/', url):
        u = urlparse(response.url)
        return u.scheme + url
    return urljoin(response.url, url)

class ConfigSpider(scrapy.Spider):
    name = 'config_spider'

    def start_requests(self):
        yield scrapy.Request(url='http://www.baidu.com/s?wd=crawlab', callback=self.parse_list)

    def parse_list(self, response):
        prev_item = response.meta.get('item')
        for elem in response.css('.result.c-container'):
            item = Item()
            item['title'] = elem.xpath('string(.//h3/a)').extract_first()
            item['url'] = elem.xpath('.//h3/a/@href').extract_first()
            item['abstract'] = elem.xpath('string(.//*[@class="c-abstract"])').extract_first()
            if prev_item is not None:
                for key, value in prev_item.items():
                    item[key] = value
            yield item
        next_url = response.css('a.n::attr("href")').extract_first()
        yield scrapy.Request(url=get_real_url(response, next_url), callback=self.parse_list, meta={'item': prev_item})


