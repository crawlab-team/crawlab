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
        yield scrapy.Request(url='https://www.amazon.cn/s?k=%E6%89%8B%E6%9C%BA&__mk_zh_CN=%E4%BA%9A%E9%A9%AC%E9%80%8A%E7%BD%91%E7%AB%99&ref=nb_sb_noss_2', callback=self.parse_list)

    def parse_list(self, response):
        prev_item = response.meta.get('item')
        for elem in response.css('.s-result-item'):
            item = Item()
            item['title'] = elem.css('span.a-text-normal::text').extract_first()
            item['url'] = elem.css('.a-link-normal::attr("href")').extract_first()
            item['price'] = elem.xpath('string(.//*[@class="a-price-whole"])').extract_first()
            item['price_fraction'] = elem.xpath('string(.//*[@class="a-price-fraction"])').extract_first()
            item['img'] = elem.css('.s-image-square-aspect > img::attr("src")').extract_first()
            if prev_item is not None:
                for key, value in prev_item.items():
                    item[key] = value
            yield item
        next_url = response.css('.a-last > a::attr("href")').extract_first()
        yield scrapy.Request(url=get_real_url(response, next_url), callback=self.parse_list, meta={'item': prev_item})


