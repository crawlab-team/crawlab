from mongoengine import *

from model.base import BaseModel


class Deploy(BaseModel):
    _id = ObjectIdField()
    spider_id = ObjectIdField()
    version = IntField()
    node_id = ObjectIdField()
