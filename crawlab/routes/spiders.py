import json
import os
import shutil
import subprocess
from datetime import datetime
from random import random

import requests
from bson import ObjectId
from flask import current_app, request
from flask_restful import reqparse, Resource
from werkzeug.datastructures import FileStorage

from config import PROJECT_DEPLOY_FILE_FOLDER, PROJECT_SOURCE_FILE_FOLDER, PROJECT_TMP_FOLDER
from constants.node import NodeStatus
from constants.task import TaskStatus
from db.manager import db_manager
from routes.base import BaseApi
from tasks.scheduler import scheduler
from tasks.spider import execute_spider
from utils import jsonify
from utils.deploy import zip_file, unzip_file
from utils.file import get_file_suffix_stats, get_file_suffix
from utils.spider import get_lang_by_stats, get_last_n_run_errors_count, get_last_n_day_tasks_count

parser = reqparse.RequestParser()
parser.add_argument('file', type=FileStorage, location='files')

IGNORE_DIRS = [
    '.idea'
]


class SpiderApi(BaseApi):
    col_name = 'spiders'

    arguments = (
        # name of spider
        ('name', str),

        # execute shell command
        ('cmd', str),

        # spider source folder
        ('src', str),

        # spider type
        ('type', str),

        # spider language
        ('lang', str),

        # spider results collection
        ('col', str),

        # spider schedule cron
        ('cron', str),

        # spider schedule cron enabled
        ('cron_enabled', int),

        # spider schedule cron enabled
        ('envs', str),

        # spider site
        ('site', str),
    )

    def get(self, id=None, action=None):
        """
        GET method of SpiderAPI.
        :param id: spider_id
        :param action: action
        """
        # action by id
        if action is not None:
            if not hasattr(self, action):
                return {
                           'status': 'ok',
                           'code': 400,
                           'error': 'action "%s" invalid' % action
                       }, 400
            return getattr(self, action)(id)

        # get one node
        elif id is not None:
            spider = db_manager.get('spiders', id=id)

            # get deploy
            last_deploy = db_manager.get_last_deploy(spider_id=spider['_id'])
            if last_deploy is not None:
                spider['deploy_ts'] = last_deploy['finish_ts']

            return jsonify(spider)

        # get a list of items
        else:
            items = []
            dirs = os.listdir(PROJECT_SOURCE_FILE_FOLDER)
            for _dir in dirs:
                if _dir in IGNORE_DIRS:
                    continue

                dir_path = os.path.join(PROJECT_SOURCE_FILE_FOLDER, _dir)
                dir_name = _dir
                spider = db_manager.get_one_by_key('spiders', key='src', value=dir_path)

                # new spider
                if spider is None:
                    stats = get_file_suffix_stats(dir_path)
                    lang = get_lang_by_stats(stats)
                    spider = db_manager.save('spiders', {
                        'name': dir_name,
                        'src': dir_path,
                        'lang': lang,
                        'suffix_stats': stats,
                    })

                # existing spider
                else:
                    # get last deploy
                    last_deploy = db_manager.get_last_deploy(spider_id=spider['_id'])
                    if last_deploy is not None:
                        spider['deploy_ts'] = last_deploy['finish_ts']

                    # get last task
                    last_task = db_manager.get_last_task(spider_id=spider['_id'])
                    if last_task is not None:
                        spider['task_ts'] = last_task['create_ts']

                    # file stats
                    stats = get_file_suffix_stats(dir_path)

                    # language
                    lang = get_lang_by_stats(stats)

                    # update spider data
                    db_manager.update_one('spiders', id=str(spider['_id']), values={
                        'lang': lang,
                        'suffix_stats': stats,
                    })

                    # ---------
                    # stats
                    # ---------
                    # last 5-run errors
                    spider['last_5_errors'] = get_last_n_run_errors_count(spider_id=spider['_id'], n=5)
                    spider['last_7d_tasks'] = get_last_n_day_tasks_count(spider_id=spider['_id'], n=5)

                # append spider
                items.append(spider)

            return {
                'status': 'ok',
                'items': jsonify(items)
            }

    def crawl(self, id: str) -> (dict, tuple):
        """
        Submit an HTTP request to start a crawl task in the node of given spider_id.
        @deprecated
        :param id: spider_id
        """
        args = self.parser.parse_args()
        node_id = args.get('node_id')

        if node_id is None:
            return {
                       'code': 400,
                       'status': 400,
                       'error': 'node_id cannot be empty'
                   }, 400

        # get node from db
        node = db_manager.get('nodes', id=node_id)

        # validate ip and port
        if node.get('ip') is None or node.get('port') is None:
            return {
                       'code': 400,
                       'status': 'ok',
                       'error': 'node ip and port should not be empty'
                   }, 400

        # dispatch crawl task
        res = requests.get('http://%s:%s/api/spiders/%s/on_crawl?node_id=%s' % (
            node.get('ip'),
            node.get('port'),
            id,
            node_id
        ))
        data = json.loads(res.content.decode('utf-8'))
        return {
            'code': res.status_code,
            'status': 'ok',
            'error': data.get('error'),
            'task': data.get('task')
        }

    def on_crawl(self, id: str) -> (dict, tuple):
        """
        Start a crawl task.
        :param id: spider_id
        :return:
        """
        args = self.parser.parse_args()
        params = args.get('params')

        spider = db_manager.get('spiders', id=ObjectId(id))

        job = execute_spider.delay(id, params)

        # create a new task
        db_manager.save('tasks', {
            '_id': job.id,
            'spider_id': ObjectId(id),
            'cmd': spider.get('cmd'),
            'params': params,
            'create_ts': datetime.utcnow(),
            'status': TaskStatus.PENDING
        })

        return {
            'code': 200,
            'status': 'ok',
            'task': {
                'id': job.id,
                'status': job.status
            }
        }

    def deploy(self, id: str) -> (dict, tuple):
        """
        Submit HTTP requests to deploy the given spider to all nodes.
        :param id:
        :return:
        """
        spider = db_manager.get('spiders', id=id)
        nodes = db_manager.list('nodes', {'status': NodeStatus.ONLINE})

        for node in nodes:
            node_id = node['_id']

            output_file_name = '%s_%s.zip' % (
                datetime.now().strftime('%Y%m%d%H%M%S'),
                str(random())[2:12]
            )
            output_file_path = os.path.join(PROJECT_TMP_FOLDER, output_file_name)

            # zip source folder to zip file
            zip_file(source_dir=spider['src'],
                     output_filename=output_file_path)

            # upload to api
            files = {'file': open(output_file_path, 'rb')}
            r = requests.post('http://%s:%s/api/spiders/%s/deploy_file?node_id=%s' % (
                node.get('ip'),
                node.get('port'),
                id,
                node_id,
            ), files=files)

            # TODO: checkpoint for errors

        return {
            'code': 200,
            'status': 'ok',
            'message': 'deploy success'
        }

    def deploy_file(self, id: str = None) -> (dict, tuple):
        """
        Receive HTTP request of deploys and unzip zip files and copy to the destination directories.
        :param id: spider_id
        """
        args = parser.parse_args()
        node_id = request.args.get('node_id')
        f = args.file

        if get_file_suffix(f.filename) != 'zip':
            return {
                       'status': 'ok',
                       'error': 'file type mismatch'
                   }, 400

        # save zip file on temp folder
        file_path = '%s/%s' % (PROJECT_TMP_FOLDER, f.filename)
        with open(file_path, 'wb') as fw:
            fw.write(f.stream.read())

        # unzip zip file
        dir_path = file_path.replace('.zip', '')
        if os.path.exists(dir_path):
            shutil.rmtree(dir_path)
        unzip_file(file_path, dir_path)

        # get spider and version
        spider = db_manager.get(col_name=self.col_name, id=id)
        if spider is None:
            return None, 400

        # make source / destination
        src = os.path.join(dir_path, os.listdir(dir_path)[0])
        # src = dir_path
        dst = os.path.join(PROJECT_DEPLOY_FILE_FOLDER, str(spider.get('_id')))

        # logging info
        current_app.logger.info('src: %s' % src)
        current_app.logger.info('dst: %s' % dst)

        # remove if the target folder exists
        if os.path.exists(dst):
            shutil.rmtree(dst)

        # copy from source to destination
        shutil.copytree(src=src, dst=dst)

        # save to db
        # TODO: task management for deployment
        db_manager.save('deploys', {
            'spider_id': ObjectId(id),
            'node_id': node_id,
            'finish_ts': datetime.utcnow()
        })

        return {
            'code': 200,
            'status': 'ok',
            'message': 'deploy success'
        }

    def get_deploys(self, id: str) -> (dict, tuple):
        """
        Get a list of latest deploys of given spider_id
        :param id: spider_id
        """
        items = db_manager.list('deploys', cond={'spider_id': ObjectId(id)}, limit=10, sort_key='finish_ts')
        deploys = []
        for item in items:
            spider_id = item['spider_id']
            spider = db_manager.get('spiders', id=str(spider_id))
            item['spider_name'] = spider['name']
            deploys.append(item)
        return {
            'status': 'ok',
            'items': jsonify(deploys)
        }

    def get_tasks(self, id: str) -> (dict, tuple):
        """
        Get a list of latest tasks of given spider_id
        :param id:
        """
        items = db_manager.list('tasks', cond={'spider_id': ObjectId(id)}, limit=10, sort_key='create_ts')
        for item in items:
            spider_id = item['spider_id']
            spider = db_manager.get('spiders', id=str(spider_id))
            item['spider_name'] = spider['name']
            if item.get('status') is None:
                item['status'] = TaskStatus.UNAVAILABLE
        return {
            'status': 'ok',
            'items': jsonify(items)
        }

    def after_update(self, id: str = None) -> None:
        """
        After each spider is updated, update the cron scheduler correspondingly.
        :param id: spider_id
        """
        scheduler.update()

    def update_envs(self, id: str):
        args = self.parser.parse_args()
        envs = json.loads(args.envs)
        db_manager.update_one(col_name='spiders', id=id, values={'envs': envs})


class SpiderImportApi(Resource):
    __doc__ = """
    API for importing spiders from external resources including Github, Gitlab, and subversion (WIP)
    """
    parser = reqparse.RequestParser()
    arguments = [
        ('url', str)
    ]

    def __init__(self):
        super(SpiderImportApi).__init__()
        for arg, type in self.arguments:
            self.parser.add_argument(arg, type=type)

    def post(self, platform: str = None) -> (dict, tuple):
        if platform is None:
            return {
                       'status': 'ok',
                       'code': 404,
                       'error': 'platform invalid'
                   }, 404

        if not hasattr(self, platform):
            return {
                       'status': 'ok',
                       'code': 400,
                       'error': 'platform "%s" invalid' % platform
                   }, 400

        return getattr(self, platform)()

    def github(self) -> None:
        """
        Import Github API
        """
        self._git()

    def gitlab(self) -> None:
        """
        Import Gitlab API
        """
        self._git()

    def _git(self):
        """
        Helper method to perform github important (basically "git clone" method).
        """
        args = self.parser.parse_args()
        url = args.get('url')
        if url is None:
            return {
                       'status': 'ok',
                       'code': 400,
                       'error': 'url should not be empty'
                   }, 400

        try:
            p = subprocess.Popen(['git', 'clone', url], cwd=PROJECT_SOURCE_FILE_FOLDER)
            _stdout, _stderr = p.communicate()
        except Exception as err:
            return {
                       'status': 'ok',
                       'code': 500,
                       'error': str(err)
                   }, 500

        return {
            'status': 'ok',
            'message': 'success'
        }


class SpiderManageApi(Resource):
    parser = reqparse.RequestParser()
    arguments = [
        ('url', str)
    ]

    def post(self, action: str) -> (dict, tuple):
        """
        POST method for SpiderManageAPI.
        :param action:
        """
        if not hasattr(self, action):
            return {
                       'status': 'ok',
                       'code': 400,
                       'error': 'action "%s" invalid' % action
                   }, 400

        return getattr(self, action)()

    def deploy_all(self) -> (dict, tuple):
        """
        Deploy all spiders to all nodes.
        """
        # active nodes
        nodes = db_manager.list('nodes', {'status': NodeStatus.ONLINE})

        # all spiders
        spiders = db_manager.list('spiders', {'cmd': {'$exists': True}})

        # iterate all nodes
        for node in nodes:
            node_id = node['_id']
            for spider in spiders:
                spider_id = spider['_id']
                spider_src = spider['src']

                output_file_name = '%s_%s.zip' % (
                    datetime.now().strftime('%Y%m%d%H%M%S'),
                    str(random())[2:12]
                )
                output_file_path = os.path.join(PROJECT_TMP_FOLDER, output_file_name)

                # zip source folder to zip file
                zip_file(source_dir=spider_src,
                         output_filename=output_file_path)

                # upload to api
                files = {'file': open(output_file_path, 'rb')}
                r = requests.post('http://%s:%s/api/spiders/%s/deploy_file?node_id=%s' % (
                    node.get('ip'),
                    node.get('port'),
                    spider_id,
                    node_id,
                ), files=files)

        return {
            'status': 'ok',
            'message': 'success'
        }
