import requests
from celery.utils.log import get_logger
from tasks import app

logger = get_logger(__name__)


@app.task
def execute_spider(spider_name: str):
    logger.info('spider_name: %s' % spider_name)
    return spider_name


@app.task
def get_baidu_html(keyword: str):
    res = requests.get('http://www.baidu.com')
    return res.content.decode('utf-8')
