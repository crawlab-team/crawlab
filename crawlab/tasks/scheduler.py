import requests
from apscheduler.schedulers.background import BackgroundScheduler
from apscheduler.jobstores.mongodb import MongoDBJobStore
from pymongo import MongoClient

from config import MONGO_DB, MONGO_HOST, MONGO_PORT, FLASK_HOST, FLASK_PORT
from constants.spider import CronEnabled
from db.manager import db_manager


class Scheduler(object):
    mongo = MongoClient(host=MONGO_HOST, port=MONGO_PORT)
    task_col = 'apscheduler_jobs'

    # scheduler jobstore
    jobstores = {
        'mongo': MongoDBJobStore(database=MONGO_DB,
                                 collection=task_col,
                                 client=mongo)
    }

    # scheduler instance
    scheduler = BackgroundScheduler(jobstores=jobstores)

    def execute_spider(self, id: str):
        r = requests.get('http://%s:%s/api/spiders/%s/on_crawl' % (
            FLASK_HOST,
            FLASK_PORT,
            id
        ))

    def update(self):
        # remove all existing periodic jobs
        self.scheduler.remove_all_jobs()
        self.mongo[MONGO_DB][self.task_col].remove()

        periodical_tasks = db_manager.list('schedules', {})
        for task in periodical_tasks:
            cron = task.get('cron')
            cron_arr = cron.split(' ')
            second = cron_arr[0]
            minute = cron_arr[1]
            hour = cron_arr[2]
            day = cron_arr[3]
            month = cron_arr[4]
            day_of_week = cron_arr[5]
            self.scheduler.add_job(func=self.execute_spider, trigger='cron', args=(str(task['spider_id']),),
                                   jobstore='mongo',
                                   day_of_week=day_of_week, month=month, day=day, hour=hour, minute=minute,
                                   second=second)

    def run(self):
        self.update()
        self.scheduler.start()


scheduler = Scheduler()

if __name__ == '__main__':
    scheduler.run()
