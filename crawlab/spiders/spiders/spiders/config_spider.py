# -*- coding: utf-8 -*-
from urllib.parse import urlparse

import scrapy

from spiders.db import spider
from spiders.items import SpidersItem


class NormalSpiderSpider(scrapy.Spider):
    name = 'config_spider'
    # allowed_domains = []
    start_urls = [spider['start_url']]

    def parse(self, response):
        if spider['item_selector_type'] == 'xpath':
            # xpath selector
            items = response.xpath(spider['item_selector'])
        else:
            # css selector
            items = response.css(spider['item_selector'])
        for _item in items:
            item = SpidersItem()
            for f in spider['fields']:
                if f['type'] == 'xpath':
                    # xpath selector
                    if f['extract_type'] == 'text':
                        # text content
                        query = f['query'] + '/text()'
                    else:
                        # attribute
                        attribute = f["attribute"]
                        query = f['query'] + f'/@("{attribute}")'
                    item[f['name']] = _item.xpath(query).extract_first()

                else:
                    # css selector
                    if f['extract_type'] == 'text':
                        # text content
                        query = f['query'] + '::text'
                    else:
                        # attribute
                        attribute = f["attribute"]
                        query = f['query'] + f'::attr("{attribute}")'
                    item[f['name']] = _item.css(query).extract_first()

                yield item

        # pagination
        if spider.get('pagination_selector') is not None:
            if spider['pagination_selector_type'] == 'xpath':
                # xpath selector
                next_url = response.xpath(spider['pagination_selector'] + '/@href').extract_first()
            else:
                # css selector
                next_url = response.css(spider['pagination_selector'] + '::attr("href")').extract_first()

            # found next url
            if next_url is not None:
                if not next_url.startswith('http') and not next_url.startswith('//'):
                    u = urlparse(response.url)
                    next_url = f'{u.scheme}://{u.netloc}{next_url}'
                yield scrapy.Request(url=next_url)
