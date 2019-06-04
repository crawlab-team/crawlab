import os
import sys
import subprocess

# make sure the working directory is in system path
FILE_DIR = os.path.dirname(os.path.realpath(__file__))
ROOT_PATH = os.path.abspath(os.path.join(FILE_DIR, '..'))
sys.path.append(ROOT_PATH)

from utils.log import other
from config import BROKER_URL

if __name__ == '__main__':
    p = subprocess.Popen([sys.executable, '-m', 'celery', 'flower', '-b', BROKER_URL],
                         stdout=subprocess.PIPE,
                         stderr=subprocess.STDOUT,
                         cwd=ROOT_PATH)
    for line in iter(p.stdout.readline, 'b'):
        if line.decode('utf-8') != '':
            other.info(line.decode('utf-8'))
