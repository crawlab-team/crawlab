from celery import Celery

# celery app instance
celery_app = Celery(__name__)
celery_app.config_from_object('config.celery')
