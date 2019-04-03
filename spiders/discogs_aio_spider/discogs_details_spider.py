# -*- coding: utf-8 -*-
# @Time : 2019/3/24 2:05 AM
# @Author : cxa
# @File : discogs_details_spider.py
# @Software: PyCharm
import asyncio

import aiofiles
from db.mongohelper import MotorOperation
from logger.log import crawler, storage

from copy import copy
import os
from common.base_crawler import Crawler
from types import AsyncGeneratorType
from decorators.decorators import decorator

from urllib.parse import urljoin
from multidict import CIMultiDict

DEFAULT_HEADRS = CIMultiDict({
    "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
    "Accept-Encoding": "gzip, deflate, br",
    "Accept-Language": "zh-CN,zh;q=0.9",
    "Host": "www.discogs.com",
    "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36",
})
BASE_URL = "https://www.discogs.com"


class Details_Spider(Crawler):
    def __init__(self):
        self.page_pat = "&page=.*&"

    @decorator()
    async def start(self):
        # 获取mongo的数据,类型异步生成器。
        data: AsyncGeneratorType = await MotorOperation().find_data(col="discogs_index_data")
        await self.init_session()
        # 分流
        tasks = (asyncio.ensure_future(self.fetch_detail_page(item)) async for item in data)
        await self.branch(tasks)

    @decorator(False)
    async def fetch_detail_page(self, item: dict):
        '''
        访问详情页，开始解析
        :param url:
        :return:

        '''
        detail_url = item.get("detail_url")
        kwargs = {"headers": DEFAULT_HEADRS}
        # 修改种子URL的状态为1表示开始爬取。
        condition = {'url': detail_url}
        await MotorOperation().change_status(condition, status_code=1)
        await asyncio.sleep(2)
        response = await self.get_session(detail_url, kwargs)
        if response.status == 200:
            source = response.source
            await self.more_images(source)
            # 获取当前的链接然后构建所有页数的url。
            # 保存当一页的内容。
            # await self.get_list_info(detail_url, source)
            # await self.max_page_index(url, source)

    @decorator(False)
    async def get_list_info(self, url, source):
        '''
        为了取得元素的正确性，这里按照块进行处理。
        :param url: 当前页的url
        :param source: 源码
        :return:
        '''
        pass
        # div_xpath = "//div[@class='cards cards_layout_text-only']/div"
        # div_node_list = self.xpath(source, div_xpath)
        # tasks = []
        # t_append = tasks.append
        # for div_node in div_node_list:
        #     try:
        #         dic = {}
        #         dic["obj_id"] = self.xpath(div_node, "@data-object-id")[0]
        #         dic["artist"] = self.xpath(div_node, ".//div[@class='card_body']/h4/span/a", "text")[0]
        #         dic["title"] = \
        #             self.xpath(div_node, ".//div[@class='card_body']/h4/a[@class='search_result_title ']", "text")[0]
        #         _detail_url = \
        #             self.xpath(div_node, ".//div[@class='card_body']/h4/a[@class='search_result_title ']", "href")[0]
        #         dic["detail_url"] = urljoin(BASE_URL, _detail_url)
        #
        #         card_info_xpath = ".//div[@class='card_body']/p[@class='card_info']"
        #         dic["label"] = self.xpath(div_node, f"{card_info_xpath}/a", "text")[0]
        #         dic["catalog_number"] = \
        #             self.xpath(div_node, f"{card_info_xpath}/span[@class='card_release_catalog_number']", "text")[0]
        #         dic["format"] = self.xpath(div_node, f"{card_info_xpath}/span[@class='card_release_format']", "text")[0]
        #         dic["year"] = self.xpath(div_node, f"{card_info_xpath}/span[@class='card_release_year']", "text")[0]
        #         dic["country"] = self.xpath(div_node, f"{card_info_xpath}/span[@class='card_release_country']", "text")[
        #             0]
        #         dic["url"] = url
        #         dic["page_index"] = 1
        #         dic["status"] = 0
        #         dic["crawler_time"] = datetime.datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        #         t_append(dic)
        #     except IndexError as e:
        #         crawler.error(f"解析出错，此时的url是:{url}")
        # await MotorOperation().save_data(dic)
        # # 修改种子URL的状态为2表示爬取成功。
        # condition = {"url": url}
        # await MotorOperation().change_status(condition, status_code=2)

    @decorator()
    async def more_images(self, source):
        '''
        获取更多图片的链接
        :param source:
        :return:
        '''
        more_url_node = self.xpath(source, "//a[contains(@class,'thumbnail_link') and contains(@href,'images')]",
                                   "href")
        if more_url_node:
            _url = more_url_node[0]
            more_url = urljoin(BASE_URL, _url)
            kwargs = {"headers": DEFAULT_HEADRS}
            response = await self.get_session(more_url, kwargs)
            if response.status == 200:
                source = response.source
                await self.parse_images(source)
                # TODO:解析歌曲的详细曲目信息

    async def get_image_buff(self, img_url):
        img_headers = copy(DEFAULT_HEADRS)
        img_headers["host"] = "img.discogs.com"
        kwargs = {"headers": img_headers}
        response = await self.get_session(img_url, kwargs, source_type="buff")
        buff = response.source
        await self.save_image(img_url, buff)

    @decorator()
    async def save_image(self, img_url, buff):
        image_name = img_url.split("/")[-1].replace(".jpeg", "")
        file_path = os.path.join(os.getcwd(), "discogs_images")
        image_path = os.path.join(file_path, image_name)
        if not os.path.exists(file_path):
            os.makedirs(file_path)
            # 文件是否存在
        if not os.path.exists(image_path):
            storage.info(f"SAVE_PATH:{image_path}")
            async with aiofiles.open(image_path, 'wb') as f:
                await f.write(buff)

    @decorator()
    async def parse_images(self, source):
        '''
        解析当前页所有图片的链接
        :param source:
        :return:
        '''
        image_node_list = self.xpath(source, "//div[@id='view_images']/p//img", "src")
        tasks = [asyncio.ensure_future(self.get_image_buff(url)) for url in image_node_list]
        await tasks


if __name__ == '__main__':
    s = Details_Spider()
    loop = asyncio.get_event_loop()
    try:
        loop.run_until_complete(s.start())
    finally:
        loop.close()
