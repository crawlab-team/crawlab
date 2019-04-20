# encoding: utf-8

import os

run_env = os.environ.get("RUNENV", "local")

if run_env == "local":  # 加载本地配置
    from config.config_local import *
else:
    from config.config import *
