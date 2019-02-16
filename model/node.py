from mongoengine import *

from model.base import BaseModel


class Node(BaseModel):
    _id = ObjectIdField()
    ip = StringField()
    port = IntField()
    name = StringField()
    description = StringField()
    status = IntField()
