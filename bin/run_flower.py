from config.celery import BROKER_URL
import subprocess

if __name__ == '__main__':
    p = subprocess.Popen(['celery', 'flower', '-b', BROKER_URL], stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    for line in iter(p.stdout.readline, 'b'):
        print(line.decode('utf-8'))
