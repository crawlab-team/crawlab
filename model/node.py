from mongoengine import *

from model.base import BaseModel


class Node(BaseModel):
    _id = ObjectIdField()
    node_ip = StringField()
    node_name = StringField()
    node_description = StringField()
