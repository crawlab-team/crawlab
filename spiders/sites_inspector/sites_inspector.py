import asyncio
import os
from datetime import datetime

import aiohttp
import requests

from pymongo import MongoClient

# MONGO_HOST = os.environ['MONGO_HOST']
# MONGO_PORT = int(os.environ['MONGO_PORT'])
# MONGO_DB = os.environ['MONGO_DB']
MONGO_HOST = 'localhost'
MONGO_PORT = 27017
MONGO_DB = 'crawlab_test'

mongo = MongoClient(host=MONGO_HOST, port=MONGO_PORT)
db = mongo[MONGO_DB]
col = db['sites']


async def process_response(resp, **kwargs):
    url = kwargs.get('url')
    status = resp.status  # 读取状态
    if status == 200 and ('robots.txt' in str(resp.url)):
        col.update({'_id': url}, {'$set': {'has_robots': True}})
    else:
        # 错误状态
        col.update({'_id': url}, {'$set': {'has_robots': False}})


async def process_home_page_response(resp, **kwargs):
    url = kwargs.get('url')
    duration = kwargs.get('duration')
    status = resp.status  # 读取状态
    col.update({'_id': url}, {'$set': {'home_http_status': status, 'home_response_time': duration}})


async def request_site(url: str, semaphore):
    _url = 'http://' + url + '/robots.txt'
    # print('crawling ' + _url)
    async with semaphore:
        async with aiohttp.ClientSession() as session:  # <1> 开启一个会话
            async with session.get(_url) as resp:  # 发送请求
                await process_response(resp=resp, url=url)
                print('crawled ' + _url)
    # resp = requests.get(_url)
    return resp


async def request_site_home_page(url: str, semophore):
    _url = 'http://' + url
    # print('crawling ' + _url)
    async with semophore:
        tic = datetime.now()
        async with aiohttp.ClientSession() as session:  # <1> 开启一个会话
            async with session.get(_url) as resp:  # 发送请求
                toc = datetime.now()
                duration = (toc - tic).total_seconds()
                await process_home_page_response(resp=resp, url=url, duration=duration)
                print('crawled ' + _url)


async def run():
    semaphore = asyncio.Semaphore(50)  # 限制并发量为50
    sites = [site for site in col.find({'rank': {'$lte': 5000}})]
    urls = [site['_id'] for site in sites]
    to_get = [request_site(url, semaphore) for url in urls]
    to_get += [request_site_home_page(url, semaphore) for url in urls]
    await asyncio.wait(to_get)


if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    loop.run_until_complete(run())
    loop.close()
