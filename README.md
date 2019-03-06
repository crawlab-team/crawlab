# Crawlab
Celery-based web crawler admin platform for managing distributed web spiders regardless of languages and frameworks.

## Pre-requisite
- Python3
- MongoDB
- Redis

## Installation

```bash
# install the requirements for backend
pip install -r ./crawlab/requirements.txt
```

```bash
cd frontend
npm install
```

## Configure

Please edit configuration file `config.py` to configure api and database connections.

## Quick Start
```bash
# run all services
python manage.py run_all
```

```bash
# run frontend client
cd frontend
npm run dev
```

## Architecture

The architecture of Crawlab is as below. It's very similar to Celery architecture, but a few more modules including Frontend, Spiders and Flower are added to feature the crawling management functionality. 

![crawlab-architecture](./docs/img/crawlab-architecture.png)

### Nodes

Nodes are actually the workers defined in Celery. A node is running and connected to a task queue, redis for example, to receive and run tasks. As spiders need to be deployed to the nodes, users should specify their ip addresses and ports before the deployment.

### Spiders

#### Auto Discovery
In `config.py` file, edit `PROJECT_SOURCE_FILE_FOLDER` as the directory where the spiders projects are located. The web app will discover spider projects automatically.

#### Deploy Spiders

All spiders need to be deployed to a specific node before crawling. Simply click "Deploy" button on spider detail page and select the right node for the deployment. 

#### Run Spiders

After deploying the spider, you can click "Run" button on spider detail page and select a specific node to start crawling. It will triggers a task for the crawling, where you can see in detail in tasks page.

### Tasks

Tasks are triggered and run by the workers. Users can check the task status info and logs in the task detail page. 
