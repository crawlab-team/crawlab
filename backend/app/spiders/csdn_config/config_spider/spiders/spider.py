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
        yield scrapy.Request(url='https://so.csdn.net/so/search/s.do?q=crawlab', callback=self.parse_list)

    def parse_list(self, response):
        prev_item = response.meta.get('item')
        for elem in response.css('.search-list-con > .search-list'):
            item = Item()
            item['url'] = elem.xpath('.//*[@class="limit_width"]/a/@href').extract_first()
            if prev_item is not None:
                for key, value in prev_item.items():
                    item[key] = value
            yield scrapy.Request(url=get_real_url(response, item['url']), callback=self.parse_detail, meta={'item': item})
        next_url = response.css('a.btn-next::attr("href")').extract_first()
        yield scrapy.Request(url=get_real_url(response, next_url), callback=self.parse_list, meta={'item': prev_item})

    def parse_detail(self, response):
        item = Item() if response.meta.get('item') is None else response.meta.get('item')
        item['content'] = response.xpath('string(.//div[@id="content_views"])').extract_first()
        item['views'] = response.css('.read-count::text').extract_first()
        item['title'] = response.css('.title-article::text').extract_first()
        item['author'] = response.css('.follow-nickName::text').extract_first()
        yield item


