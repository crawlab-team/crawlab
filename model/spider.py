from mongoengine import *

from model.base import BaseModel


class Spider(BaseModel):
    _id = ObjectIdField()
    spider_name = StringField()
    spider_type = IntField()
    lang_type = IntField()
    execute_cmd = StringField()
    src_file_path = StringField()
