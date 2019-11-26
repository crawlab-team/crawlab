# -*- coding: utf-8 -*-
import scrapy
import re
from config_spider.items import Item
from urllib.parse import urljoin

def get_real_url(response, url):
    if re.search(r'^https?|^\/\/', url):
        return url
    return urljoin(response.url, url)

class ConfigSpider(scrapy.Spider):
    name = 'config_spider'

    def start_requests(self):
        yield scrapy.Request(url='###START_URL###', callback=self.###START_STAGE###)

###PARSERS###
