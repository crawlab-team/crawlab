# -*- coding: utf-8 -*-
import scrapy
from chinaz.items import ChinazItem


class ChinazSpiderSpider(scrapy.Spider):
    name = 'chinaz_spider'
    allowed_domains = ['chinaz.com']
    start_urls = ['http://top.chinaz.com/hangye/']

    def parse(self, response):
        for item in response.css('.listCentent > li'):
            name = item.css('h3.rightTxtHead > a::text').extract_first()
            href = item.css('h3.rightTxtHead > a::attr("href")').extract_first()
            domain = item.css('h3.rightTxtHead > span::text').extract_first()
            description = item.css('p.RtCInfo::text').extract_first()
            rank = item.css('.RtCRateCent > strong::text').extract_first()
            rank = int(rank)
            item = ChinazItem(
                _id=domain,
                name=name,
                domain=domain,
                description=description,
                rank=rank,
            )
            yield scrapy.Request(
                url='http://top.chinaz.com' + href,
                callback=self.parse_item,
                meta={
                    'item': item
                }
            )

        # pagination
        a_list = response.css('.ListPageWrap > a::attr("href")').extract()
        url = 'http://top.chinaz.com/hangye/' + a_list[-1]
        yield scrapy.Request(url=url, callback=self.parse)

    def parse_item(self, response):
        item = response.meta['item']

        # category info extraction
        arr = response.css('.TopMainTag-show .SimSun')
        res1 = arr[0].css('a::text').extract()
        main_category = res1[0]
        if len(res1) == 1:
            category = '其他'
        else:
            category = res1[1]

        # location info extraction
        res2 = arr[1].css('a::text').extract()
        if len(res2) > 0:
            location = res2[0]
        else:
            location = '其他'

        # assign values to item
        item['main_category'] = main_category
        item['category'] = category
        item['location'] = location

        yield item
