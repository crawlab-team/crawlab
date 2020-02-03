import requests
from bs4 import BeautifulSoup as bs
from urllib.parse import urljoin, urlparse
import re
from crawlab import save_item

s = requests.Session()

def get_real_url(response, url):
    if re.search(r'^https?', url):
        return url
    elif re.search(r'^\/\/', url):
        u = urlparse(response.url)
        return u.scheme + url
    return urljoin(response.url, url)

def start_requests():
	for i in range(0, 9):
		fr = 'PERE' if not i else 'MORE'
		url = f'https://cn.bing.com/search?q=crawlab&first={10 * i + 1}&FROM={fr}'
		request_page(url)

def request_page(url):
	print(f'requesting {url}')
	r = s.get(url, headers={'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36'})
	parse_list(r)

def parse_list(response):
	soup = bs(response.content.decode('utf-8'))
	for el in list(soup.select('#b_results > li')):
		try:
			save_item({
				'title': el.select_one('h2').text,
				'url': el.select_one('h2 a').attrs.get('href'),
				'abstract': el.select_one('.b_caption p').text,
			})
		except:
			pass

if __name__ == '__main__':
	start_requests()