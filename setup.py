#-*- encoding: UTF-8 -*-
from setuptools import setup, find_packages

VERSION = '0.2.3'

with open('README.md') as fp:
    readme = fp.read()

setup(name='crawlab-server',
      version=VERSION,
      description="Celery-based web crawler admin platform for managing distributed web spiders regardless of languages and frameworks.",
      long_description=readme,
      classifiers=['Python', 'Javascript', 'Scrapy'], # Get strings from http://pypi.python.org/pypi?%3Aaction=list_classifiers
      keywords='python crawlab celery crawler spider platform scrapy',
      author='tikazyq',
      author_email='tikazyq@163.com',
      url='https://github.com/tikazyq/crawlab',
      license='BSD',
      packages=find_packages(),
      include_package_data=True,
      zip_safe=True,
      install_requires=[
        'celery',
        'flower',
        'requests',
        'pymongo',
        'flask',
        'flask_cors',
        'flask_restful',
        'lxml',
        'gevent',
        'scrapy',
      ],
      entry_points={
        'console_scripts':[
            'crawlab = crawlab.manage:main'
        ]
      },
)