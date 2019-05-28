# -*- coding: utf-8 -*-
import os
import sys
from urllib.parse import urlparse

import scrapy

from spiders.db import spider
from spiders.items import SpidersItem
from spiders.utils import generate_urls

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..', '..')))


def get_detail_url(item):
    for f in spider['fields']:
        if f.get('is_detail'):
            return item.get(f['name'])
    return None


def get_spiders_item(sel, fields, item=None):
    if item is None:
        item = SpidersItem()

    for f in fields:
        if f['type'] == 'xpath':
            # xpath selector
            if f['extract_type'] == 'text':
                # text content
                query = f['query'] + '/text()'
            else:
                # attribute
                attribute = f["attribute"]
                query = f['query'] + f'/@("{attribute}")'
            item[f['name']] = sel.xpath(query).extract_first()

        else:
            # css selector
            if f['extract_type'] == 'text':
                # text content
                query = f['query'] + '::text'
            else:
                # attribute
                attribute = f["attribute"]
                query = f['query'] + f'::attr("{attribute}")'
            item[f['name']] = sel.css(query).extract_first()

    return item


def get_list_items(response):
    if spider['item_selector_type'] == 'xpath':
        # xpath selector
        items = response.xpath(spider['item_selector'])
    else:
        # css selector
        items = response.css(spider['item_selector'])
    return items


def get_next_url(response):
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
                return next_url
    return None


class ConfigSpiderSpider(scrapy.Spider):
    name = 'config_spider'

    def start_requests(self):
        for url in generate_urls(spider['start_url']):
            yield scrapy.Request(url=url)

    def parse(self, response):

        if spider['crawl_type'] == 'list':
            # list page only
            items = get_list_items(response)
            for _item in items:
                item = get_spiders_item(sel=_item, fields=spider['fields'])
                yield item
            next_url = get_next_url(response)
            if next_url is not None:
                yield scrapy.Request(url=next_url)

        elif spider['crawl_type'] == 'detail':
            # TODO: detail page only
            # detail page only
            pass

        elif spider['crawl_type'] == 'list-detail':
            # list page + detail page
            items = get_list_items(response)
            for _item in items:
                item = get_spiders_item(sel=_item, fields=spider['fields'])
                detail_url = get_detail_url(item)
                if detail_url is not None:
                    yield scrapy.Request(url=detail_url,
                                         callback=self.parse_detail,
                                         meta={
                                             'item': item
                                         })
            next_url = get_next_url(response)
            if next_url is not None:
                yield scrapy.Request(url=next_url)

    def parse_detail(self, response):
        item = get_spiders_item(sel=response, fields=spider['detail_fields'], item=response.meta['item'])
        yield item
