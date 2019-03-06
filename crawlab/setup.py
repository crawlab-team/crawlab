from setuptools import setup, find_packages

with open("README.md", "r") as fh:
    long_description = fh.read()

with open('requirements.txt') as f:
    requirements = [l for l in f.read().splitlines() if l]

setup(
    name='crawlab-server',
    version='0.0.1',
    url='https://github.com/tikazyq/crawlab',
    install_requires=requirements,
    license='BSD',
    author='Marvin Zhang',
    author_email='tikazyq@163.com',
    description='Celery-based web crawler admin platform for managing distributed web spiders regardless of languages and frameworks.',
    long_description=long_description,
    long_description_content_type="text/markdown",
    download_url="https://github.com/tikazyq/crawlab/archive/master.zip",
    packages=find_packages(),
    keywords=['celery', 'python', 'webcrawler', 'crawl', 'scrapy', 'admin'],
    zip_safe=True,
)
