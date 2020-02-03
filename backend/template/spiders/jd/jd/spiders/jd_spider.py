# -*- coding: utf-8 -*-
import scrapy

from jd.items import JdItem


class JdSpiderSpider(scrapy.Spider):
    name = 'jd_spider'
    allowed_domains = ['jd.com']

    def start_requests(self):
    	for i in range(1, 50):
    		yield scrapy.Request(url=f'https://search.jd.com/Search?keyword=手机&enc=utf-8&page={i}')

    def parse(self, response):
        for el in response.css('.gl-item'):
            yield JdItem(
                url=el.css('.p-name > a::attr("href")').extract_first(),
                name=el.css('.p-name > a::attr("title")').extract_first(),
                price=float(el.css('.p-price i::text').extract_first()),
            )
