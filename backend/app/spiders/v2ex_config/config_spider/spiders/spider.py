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
        yield scrapy.Request(url='https://v2ex.com/', callback=self.parse_list)

    def parse_list(self, response):
        prev_item = response.meta.get('item')
        for elem in response.css('.cell.item'):
            item = Item()
            item['title'] = elem.css('a.topic-link::text').extract_first()
            item['url'] = elem.css('a.topic-link::attr("href")').extract_first()
            item['replies'] = elem.css('.count_livid::text').extract_first()
            if prev_item is not None:
                for key, value in prev_item.items():
                    item[key] = value
            yield scrapy.Request(url=get_real_url(response, item['url']), callback=self.parse_detail, meta={'item': item})

    def parse_detail(self, response):
        item = Item() if response.meta.get('item') is None else response.meta.get('item')
        item['content'] = response.xpath('string(.//*[@class="markdown_body"])').extract_first()
        yield item


