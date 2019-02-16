from mongoengine import *

from model.base import BaseModel


class Spider(BaseModel):
    _id = ObjectIdField()
    name = StringField()
    cmd = StringField()
    src = StringField()
    type = IntField()
    lang = IntField()
