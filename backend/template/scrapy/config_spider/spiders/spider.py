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
        yield scrapy.Request(url='###START_URL###', callback=self.###START_STAGE###)

###PARSERS###
