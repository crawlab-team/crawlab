# -*- coding: utf-8 -*-
# @Time : 2019/3/22 10:02 PM
# @Author : cxa
# @File : testdic.py
# @Software: PyCharm
from copy import deepcopy
img_url="http:"
try:
   a=img_url.split("/")[-1]
   print(a)
   a.replace(".jpeg","")
except Exception as e:
    print(e.args)
