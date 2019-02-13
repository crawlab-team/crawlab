from mongoengine import *

from model.base import BaseModel


class Task(BaseModel):
    _id = ObjectIdField()
    deploy_id = ObjectIdField()
    file_path = StringField()
