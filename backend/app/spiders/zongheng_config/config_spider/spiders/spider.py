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
        yield scrapy.Request(url='http://www.zongheng.com/rank/details.html?rt=1&d=1', callback=self.parse_list)

    def parse_list(self, response):
        prev_item = response.meta.get('item')
        for elem in response.css('.rank_d_list'):
            item = Item()
            item['title'] = elem.css('.rank_d_b_name > a::text').extract_first()
            item['url'] = elem.css('.rank_d_b_name > a::attr("href")').extract_first()
            item['abstract'] = elem.css('body::text').extract_first()
            item['votes'] = elem.css('.rank_d_b_ticket::text').extract_first()
            if prev_item is not None:
                for key, value in prev_item.items():
                    item[key] = value
            yield item


