# -*- coding: utf-8 -*-
# @Time : 2019-03-22 15:05
# @Author : cxa
# @File : base_crawler.py
# @Software: PyCharm

import asyncio
import aiohttp
from logger.log import crawler
import async_timeout
from collections import namedtuple
from config.config import *
from multidict import CIMultiDict
from typing import Optional, Union
from async_retrying import retry
from lxml import html
from aiostream import stream
import marshal

Response = namedtuple("Response",
                      ["status", "source"])

try:
    import uvloop

    asyncio.set_event_loop_policy(uvloop.EventLoopPolicy())
except ImportError:
    pass
sem = asyncio.Semaphore(CONCURRENCY_NUM)


class Crawler():
    def __init__(self):
        self.session = None
        self.tc = None

    @retry(attempts=MAX_RETRY_TIMES)
    async def get_session(self, url, _kwargs: dict = {}, source_type="text", status_code=200) -> Response:
        """
        :param kwargs:url,headers,data,params,etc,,
        :param method: get post.
        :param timeout: defalut 5s.
        """
        kwargs = marshal.loads(marshal.dumps(_kwargs))
        if USE_PROXY:
            kwargs["proxy"] = await self.get_proxy()
        method = kwargs.pop("method", "get")
        timeout = kwargs.pop("timeout", 5)
        with async_timeout.timeout(timeout):
            async with getattr(self.session, method)(url, **kwargs) as req:
                status = req.status
                if status in [status_code, 201]:
                    if source_type == "text":
                        source = await req.text()
                    elif source_type == "buff":
                        source = await req.read()

        crawler.info(f"get url:{url},status:{status}")
        res = Response(status=status, source=source)
        return res

    def xpath(self, _response, rule, _attr=None):
        if isinstance(_response, Response):
            source = _response.text
            root = html.fromstring(source)

        elif isinstance(_response, str):
            source = _response
            root = html.fromstring(source)
        else:
            root = _response
        nodes = root.xpath(rule)
        result = []
        if _attr:
            if _attr == "text":
                result = (entry.text for entry in nodes)
            else:
                result = (entry.get(_attr) for entry in nodes)
        else:
            result = nodes
        return result

    async def branch(self, coros, limit=10):
        """
        使用aiostream模块对异步生成器做一个切片操作。这里并发量为10.
        :param coros: 异步生成器
        :param limit: 并发次数
        :return:
        """
        index = 0
        while True:
            xs = stream.preserve(coros)
            ys = xs[index:index + limit]
            t = await stream.list(ys)
            if not t:
                break
            await asyncio.ensure_future(asyncio.wait(t))
            index += limit + 1
    def call_back(self):
        return "请输入get或者post"
    async def get_proxy(self) -> Optional[str]:
        ...

    async def init_session(self, cookies=None):
        """
        创建Tcpconnector，包括ssl和连接数的限制
        创建一个全局session。
        :return:
        """
        self.tc = aiohttp.connector.TCPConnector(limit=300, force_close=True,
                                                 enable_cleanup_closed=True,
                                                 verify_ssl=False)
        self.session = aiohttp.ClientSession(connector=self.tc, cookies=cookies)

    def run(self):
        '''
        创建全局session
        :return:
        '''
        loop = asyncio.get_event_loop()
        loop.run_until_complete(self.get_session("url",{"method":"poo"}))

    async def close(self):
        await self.tc.close()
        await self.session.close()


if __name__ == '__main__':
    c = Crawler().run()
