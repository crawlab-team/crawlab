# -*- coding: utf-8 -*-
import json
from time import sleep

import scrapy

from xueqiu.items import XueqiuItem


class XueqiuSpiderSpider(scrapy.Spider):
    name = 'xueqiu_spider'
    allowed_domains = ['xueqiu.com']

    def start_requests(self):
        return [scrapy.Request(
            url='https://xueqiu.com',
            callback=self.parse_home
        )]

    def parse_home(self, response):
        yield scrapy.Request(
            url='https://xueqiu.com/v4/statuses/public_timeline_by_category.json?since_id=-1&max_id=-1&count=20&category=6'
        )

    def parse(self, response):
        data = json.loads(response.body)
        next_max_id = data.get('next_max_id')
        sleep(1)
        for row in data.get('list'):
            d = json.loads(row.get('data'))
            item = XueqiuItem(
                id=d['id'],
                text=d['text'],
                mark=d['mark'],
                target=d['target'],
                created_at=d['created_at'],
                view_count=d['view_count'],
            )
            yield item

        yield scrapy.Request(
            url=f'https://xueqiu.com/v4/statuses/public_timeline_by_category.json?since_id=-1&max_id={next_max_id}&count=20&category=6'
        )
