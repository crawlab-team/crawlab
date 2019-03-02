import json
import os
import shutil
from datetime import datetime
from random import random

import requests
from bson import ObjectId
from flask import current_app, request
from flask_restful import reqparse
from werkzeug.datastructures import FileStorage

from config import PROJECT_DEPLOY_FILE_FOLDER, PROJECT_SOURCE_FILE_FOLDER, PROJECT_TMP_FOLDER
from db.manager import db_manager
from routes.base import BaseApi
from tasks.spider import execute_spider
from utils import jsonify
from utils.deploy import zip_file, unzip_file
from utils.file import get_file_suffix_stats, get_file_suffix
from utils.spider import get_lang_by_stats

parser = reqparse.RequestParser()
parser.add_argument('file', type=FileStorage, location='files')


class SpiderApi(BaseApi):
    col_name = 'spiders'

    arguments = (
        ('name', str),
        ('cmd', str),
        ('src', str),
        ('type', str),
        ('lang', str),

        # for deploy only
        ('node_id', str),
    )

    def get(self, id=None, action=None):
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
            return jsonify(db_manager.get('spiders', id=id))

        # get a list of items
        else:
            dirs = os.listdir(PROJECT_SOURCE_FILE_FOLDER)
            for _dir in dirs:
                dir_path = os.path.join(PROJECT_SOURCE_FILE_FOLDER, _dir)
                dir_name = _dir
                spider = db_manager.get_one_by_key('spiders', key='src', value=dir_path)

                # new spider
                if spider is None:
                    stats = get_file_suffix_stats(dir_path)
                    lang = get_lang_by_stats(stats)
                    db_manager.save('spiders', {
                        'name': dir_name,
                        'src': dir_path,
                        'lang': lang,
                        'suffix_stats': stats,
                    })

                # existing spider
                else:
                    stats = get_file_suffix_stats(dir_path)
                    lang = get_lang_by_stats(stats)
                    db_manager.update_one('spiders', id=str(spider['_id']), values={
                        'lang': lang,
                        'suffix_stats': stats,
                    })

            items = db_manager.list('spiders', {})
            for item in items:
                last_deploy = db_manager.get_last_deploy(spider_id=str(item['_id']))
                if last_deploy:
                    item['update_ts'] = last_deploy['finish_ts'].strftime('%Y-%m-%d %H:%M:%S')

        return jsonify({
            'status': 'ok',
            'items': items
        })

    def crawl(self, id):
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

    def on_crawl(self, id):
        args = self.parser.parse_args()
        node_id = args.get('node_id')

        job = execute_spider.delay(id, node_id)

        return {
            'code': 200,
            'status': 'ok',
            'task': {
                'id': job.id,
                'status': job.status
            }
        }

    def deploy(self, id):
        args = self.parser.parse_args()
        node_id = args.get('node_id')

        # get spider given the id
        spider = db_manager.get(col_name=self.col_name, id=id)
        if spider is None:
            return

        # get node given the node
        node = db_manager.get(col_name='nodes', id=node_id)

        # get latest version
        latest_version = db_manager.get_latest_version(spider_id=id, node_id=node_id)

        # initialize version if no version found
        if latest_version is None:
            latest_version = 0

        # make source / destination
        src = spider.get('src')

        # copy files
        # TODO: multi-node copy files
        try:
            # TODO: deploy spiders to other node rather than in local machine
            output_file_name = '%s_%s.zip' % (
                datetime.now().strftime('%Y%m%d%H%M%S'),
                str(random())[2:12]
            )
            output_file_path = os.path.join(PROJECT_TMP_FOLDER, output_file_name)

            # zip source folder to zip file
            zip_file(source_dir=src,
                     output_filename=output_file_path)

            # upload to api
            files = {'file': open(output_file_path, 'rb')}
            r = requests.post('http://%s:%s/api/spiders/%s/deploy_file?node_id=%s' % (
                node.get('ip'),
                node.get('port'),
                id,
                node_id,
            ), files=files)

            if r.status_code == 200:
                return {
                    'code': 200,
                    'status': 'ok',
                    'message': 'deploy success'
                }

            else:
                return {
                           'code': r.status_code,
                           'status': 'ok',
                           'error': r.content.decode('utf-8')
                       }, r.status_code

        except Exception as err:
            current_app.logger.error(err)
            return {
                       'code': 500,
                       'status': 'ok',
                       'error': str(err)
                   }, 500

        finally:
            version = latest_version + 1
            db_manager.save('deploys', {
                'spider_id': ObjectId(id),
                'version': version,
                'node_id': node_id,
                'finish_ts': datetime.now()
            })

    def deploy_file(self, id=None):
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

        # get version
        latest_version = db_manager.get_latest_version(spider_id=id, node_id=node_id)
        if latest_version is None:
            latest_version = 0

        # make source / destination
        src = os.path.join(dir_path, os.listdir(dir_path)[0])
        # src = dir_path
        dst = os.path.join(PROJECT_DEPLOY_FILE_FOLDER, str(spider.get('_id')), str(latest_version + 1))

        # logging info
        current_app.logger.info('src: %s' % src)
        current_app.logger.info('dst: %s' % dst)

        # copy from source to destination
        shutil.copytree(src=src, dst=dst)

        return {
            'code': 200,
            'status': 'ok',
            'message': 'deploy success'
        }

    def get_deploys(self, id):
        items = db_manager.list('deploys', cond={'spider_id': ObjectId(id)}, limit=10, sort_key='finish_ts')
        deploys = []
        for item in items:
            spider_id = item['spider_id']
            spider = db_manager.get('spiders', id=str(spider_id))
            item['spider_name'] = spider['name']
            deploys.append(item)
        return jsonify({
            'status': 'ok',
            'items': deploys
        })

    def get_tasks(self, id):
        items = db_manager.list('tasks', cond={'spider_id': ObjectId(id)}, limit=10, sort_key='finish_ts')
        for item in items:
            spider_id = item['spider_id']
            spider = db_manager.get('spiders', id=str(spider_id))
            item['spider_name'] = spider['name']
            task = db_manager.get('tasks_celery', id=item['_id'])
            if task is not None:
                item['status'] = task['status']
            else:
                item['status'] = 'UNAVAILABLE'
        return jsonify({
            'status': 'ok',
            'items': items
        })
