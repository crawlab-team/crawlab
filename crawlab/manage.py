import argparse
import os
import subprocess
from multiprocessing import Process
import sys

BASE_DIR = os.path.dirname(__file__)

APP_DESC = """
Crawlab CLI tool. 
"""
ACTION_LIST = [
    'serve',
    'app',
    'worker',
    'flower',
    'frontend',
]
if len(sys.argv) == 1:
    sys.argv.append('--help')
parser = argparse.ArgumentParser()
parser.add_argument('action', type=str)
# parser.add_argument('-q', '--quality', type=int, default=0,
#                     help="download video quality : 1 for the standard-definition; 3 for the super-definition")
args = parser.parse_args()


def run_app():
    p = subprocess.Popen([sys.executable, os.path.join(BASE_DIR, 'app.py')])
    p.communicate()


def run_flower():
    p = subprocess.Popen([sys.executable, os.path.join(BASE_DIR, 'flower.py')])
    p.communicate()


def run_worker():
    p = subprocess.Popen([sys.executable, os.path.join(BASE_DIR, 'worker.py')])
    p.communicate()


def run_frontend():
    p = subprocess.Popen(['npm', 'run', 'serve'],
                         cwd=os.path.abspath(os.path.join(BASE_DIR, '..', 'frontend')))
    p.communicate()


def main():
    p_app = Process(target=run_app)
    p_flower = Process(target=run_flower)
    p_worker = Process(target=run_worker)
    p_frontend = Process(target=run_frontend)
    if args.action == 'serve':
        p_app.start()
        p_flower.start()
        p_worker.start()
        p_frontend.start()
    elif args.action == 'app':
        p_app.start()
        p_flower.start()
    elif args.action == 'worker':
        p_app.start()
        p_worker.start()
    elif args.action == 'flower':
        p_flower.start()
    elif args.action == 'frontend':
        p_frontend.start()
    else:
        print(f'Invalid action: {args.action}')


if __name__ == '__main__':
    main()
