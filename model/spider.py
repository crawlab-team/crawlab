from mongoengine import *

from model.base import BaseModel


class Spider(BaseModel):
    _id = ObjectIdField()
    spider_name = StringField()
    cmd = StringField()
    src = StringField()
    spider_type = IntField()
    lang_type = IntField()
