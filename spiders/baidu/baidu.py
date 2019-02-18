from time import sleep
import requests

r = requests.get('http://www.baidu.com')
print(r.content)
