import os
import shutil
from datetime import datetime

from bson import ObjectId

from config.common import PROJECT_DEPLOY_FILE_FOLDER, PROJECT_SOURCE_FILE_FOLDER
from db.manager import db_manager
from routes.base import BaseApi
from tasks.spider import execute_spider
from utils import jsonify
from utils.file import get_file_suffix_stats
from utils.spider import get_lang_by_stats


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

        return jsonify({
            'status': 'ok',
            'items': items
        })

    def crawl(self, id):
        args = self.parser.parse_args()
        node_id = args.get('node_id')

        if node_id is None:
            return {}, 400

        job = execute_spider.delay(id, node_id)
        # print('crawl: %s' % id)
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

        # TODO: deploy spiders to other node rather than in local machine

        # get latest version
        latest_version = db_manager.get_latest_version(spider_id=id)

        # initialize version if no version found
        if latest_version is None:
            latest_version = 0

        # make source / destination
        src = spider.get('src')
        dst = os.path.join(PROJECT_DEPLOY_FILE_FOLDER, str(spider.get('_id')), str(latest_version + 1))

        # copy files
        try:
            shutil.copytree(src=src, dst=dst)
            return {
                'code': 200,
                'status': 'ok',
                'message': 'deploy success'
            }
        except Exception as err:
            print(err)
            return {
                'code': 500,
                'status': 'ok',
                'error': str(err)
            }
        finally:
            version = latest_version + 1
            db_manager.save('deploys', {
                'spider_id': ObjectId(id),
                'version': version,
                'node_id': node_id,
                'finish_ts': datetime.now()
            })

    def get_deploys(self, id):
        items = db_manager.list('deploys', {'spider_id': ObjectId(id)}, limit=10)
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
        items = db_manager.list('tasks', {'spider_id': ObjectId(id)}, limit=10)
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
