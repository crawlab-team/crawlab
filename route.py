from app import api
from api.spider import SpiderApi, SpiderExecutorApi

api.add_resource(SpiderExecutorApi, '/spider')
