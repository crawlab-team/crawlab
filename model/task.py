from mongoengine import *

from model.base import BaseModel


class Task(BaseModel):
    node_id = ObjectIdField()
    spider_id = ObjectIdField()
    file_path = StringField()
