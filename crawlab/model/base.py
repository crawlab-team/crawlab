from mongoengine import *
import datetime


class BaseModel(Document):
    create_ts = DateTimeField(default=datetime.datetime.utcnow)
