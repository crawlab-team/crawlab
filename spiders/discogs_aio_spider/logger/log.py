import os
import logging
import logging.config as log_conf
import datetime
import coloredlogs

log_dir = os.path.dirname(os.path.dirname(__file__)) + '/logs'
if not os.path.exists(log_dir):
    os.mkdir(log_dir)
today = datetime.datetime.now().strftime("%Y%m%d")

log_path = os.path.join(log_dir, f'discogs{today}.log')

log_config = {
    'version': 1.0,
    'formatters': {
        'colored_console': {'()': 'coloredlogs.ColoredFormatter',
                            'format': "%(asctime)s - %(name)s - %(levelname)s - %(message)s", 'datefmt': '%H:%M:%S'},
        'detail': {
            'format': '%(asctime)s - %(name)s - %(levelname)s - %(message)s',
            'datefmt': "%Y-%m-%d %H:%M:%S"  # 如果不加这个会显示到毫秒。
        },
        'simple': {
            'format': '%(name)s - %(levelname)s - %(message)s',
        },
    },
    'handlers': {
        'console': {
            'class': 'logging.StreamHandler',  # 日志打印到屏幕显示的类。
            'level': 'INFO',
            'formatter': 'colored_console'
        },
        'file': {
            'class': 'logging.handlers.RotatingFileHandler',  # 日志打印到文件的类。
            'maxBytes': 1024 * 1024 * 1024,  # 单个文件最大内存
            'backupCount': 1,  # 备份的文件个数
            'filename': log_path,  # 日志文件名
            'level': 'INFO',  # 日志等级
            'formatter': 'detail',  # 调用上面的哪个格式
            'encoding': 'utf-8',  # 编码
        },
    },
    'loggers': {
        'crawler': {
            'handlers': ['console', 'file'],  # 只打印屏幕
            'level': 'DEBUG',  # 只显示错误的log
        },
        'parser': {
            'handlers': ['file'],
            'level': 'INFO',
        },
        'other': {
            'handlers': ['console', 'file'],
            'level': 'INFO',
        },
        'storage': {
            'handlers': ['console'],
            'level': 'INFO',
        }
    }
}

log_conf.dictConfig(log_config)

crawler = logging.getLogger('crawler')
storage = logging.getLogger('storage')
coloredlogs.install(level='DEBUG', logger=crawler)
coloredlogs.install(level='DEBUG', logger=storage)
