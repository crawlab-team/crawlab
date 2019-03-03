from setuptools import setup

with open("README.md", "r") as fh:
    long_description = fh.read()

with open('requirements.txt') as f:
    requirements = [l for l in f.read().splitlines() if l]

setup(
    name='crawlab',
    version='0.0.1',
    packages=['db', 'test', 'model', 'tasks', 'utils', 'routes', 'constants'],
    url='https://github.com/tikazyq/crawlab',
    license='BSD',
    author='Marvin Zhang',
    author_email='tikazyq@163.com',
    description='Celery-based web crawler admin platform for managing distributed web spiders regardless of languages and frameworks.'
)
