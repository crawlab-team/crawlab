# -*- coding: utf-8 -*-
import scrapy

from realestate.items import RealEstateItem


class LianjiaSpider(scrapy.Spider):
    name = 'lianjia'
    allowed_domains = ['lianjia.com']
    start_urls = ['https://cq.lianjia.com/ershoufang/']

    def start_requests(self):
        for i in range(100):
            url = 'https://cq.lianjia.com/ershoufang/pg%s' % i
            yield scrapy.Request(url=url)

    def parse(self, response):
        for item in response.css('.sellListContent > li'):
            yield RealEstateItem(
                name=item.css('.title > a::text').extract_first(),
                url=item.css('.title > a::attr("href")').extract_first(),
                type='secondhand',
                price=item.css('.totalPrice > span::text').extract_first(),
                region=item.css('.houseInfo > a::text').extract_first(),
                size=item.css('.houseInfo::text').extract_first().split(' | ')[2]
            )

        # 分页
        # a_next = response.css('.house-lst-page-box > a')[-1]
        # href = a_next.css('a::attr("href")')
        # yield scrapy.Response(url='https://cq.lianjia.com' + href)
