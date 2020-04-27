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
        yield scrapy.Request(url='https://book.douban.com/latest', callback=self.parse_list)

    def parse_list(self, response):
        prev_item = response.meta.get('item')
        for elem in response.css('ul.cover-col-4 > li'):
            item = Item()
            item['title'] = elem.css('h2 > a::text').extract_first()
            item['url'] = elem.css('h2 > a::attr("href")').extract_first()
            item['img'] = elem.css('a.cover img::attr("src")').extract_first()
            item['rating'] = elem.css('p.rating > .color-lightgray::text').extract_first()
            item['abstract'] = elem.css('p:last-child::text').extract_first()
            item['info'] = elem.css('.color-gray::text').extract_first()
            if prev_item is not None:
                for key, value in prev_item.items():
                    item[key] = value
            yield item


