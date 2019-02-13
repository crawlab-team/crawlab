from time import sleep
import requests

for i in range(10):
    r = requests.get('http://www.baidu.com')
    sleep(0.1)
