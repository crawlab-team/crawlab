# -*- coding: utf-8 -*-
import os
import re
from datetime import datetime

import scrapy
from pymongo import MongoClient

from sinastock.items import NewsItem

class SinastockSpiderSpider(scrapy.Spider):
    name = 'sinastock_spider'
    allowed_domains = ['finance.sina.com.cn']

    def start_requests(self):
        col = self.db['stocks']
        for s in col.find({}):
            code, ex = s['ts_code'].split('.')
            for i in range(10):
                url = f'http://vip.stock.finance.sina.com.cn/corp/view/vCB_AllNewsStock.php?symbol={ex.lower()}{code}&Page={i + 1}'
                yield scrapy.Request(
                    url=url,
                    callback=self.parse,
                    meta={'ts_code': s['ts_code']}
                )

    def parse(self, response):
        for a in response.css('.datelist > ul > a'):
            url = a.css('a::attr("href")').extract_first()
            item = NewsItem(
                title=a.css('a::text').extract_first(),
                url=url,
                source='sina',
                stocks=[response.meta['ts_code']]
            )
            yield scrapy.Request(
                url=url,
                callback=self.parse_detail,
                meta={'item': item}
            )

    def parse_detail(self, response):
        item = response.meta['item']
        text = response.css('#artibody').extract_first()
        pre = re.compile('>(.*?)<')
        text = ''.join(pre.findall(text))
        item['text'] = text.replace('\u3000', '')
        item['ts_str'] = response.css('.date::text').extract_first()
        if item['text'] is None or item['ts_str'] is None:
            pass
        else:
            item['ts'] = datetime.strptime(item['ts_str'], '%Y年%m月%d日 %H:%M')
            yield item
