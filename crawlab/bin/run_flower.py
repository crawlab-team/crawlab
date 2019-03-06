import os
import sys
import subprocess

# make sure the working directory is in system path
file_dir = os.path.dirname(os.path.realpath(__file__))
root_path = os.path.abspath(os.path.join(file_dir, '..'))
sys.path.append(root_path)

from config import BROKER_URL

if __name__ == '__main__':
    p = subprocess.Popen(['celery', 'flower', '-b', BROKER_URL], stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    for line in iter(p.stdout.readline, 'b'):
        if line.decode('utf-8') != '':
            print(line.decode('utf-8'))
