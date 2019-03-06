import sys
import os

# make sure the working directory is in system path
file_dir = os.path.dirname(os.path.realpath(__file__))
root_path = os.path.abspath(os.path.join(file_dir, '..'))
sys.path.append(root_path)

from config import PROJECT_LOGS_FOLDER, FLASK_HOST, FLASK_PORT
from manage import app

# create folder if it does not exist
if not os.path.exists(PROJECT_LOGS_FOLDER):
    os.makedirs(PROJECT_LOGS_FOLDER)

# run app instance
app.run(host=FLASK_HOST, port=FLASK_PORT)
