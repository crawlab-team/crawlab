pipeline {
    agent {
        node {
            label 'crawlab'
        }
    }

    environment {
        NODE_HOME = '/home/yeqing/.nvm/versions/node/v8.12.0'
        ROOT_DIR = "/home/yeqing/jenkins_home/workspace/crawlab_${GIT_BRANCH}" 
    }

    stages {
        stage('Setup') {
            steps {
                echo "Running Setup..."

                sh '#source /home/yeqing/.profile'
            }
        }
        stage('Build Frontend') {
            steps {
                echo "Building frontend..."
                sh "cd ${ROOT_DIR}/crawlab && ${NODE_HOME}/node ${NODE_HOME}/bin/npm install"
                sh "cd ${ROOT_DIR}/crawlab && ${NODE_HOME}/bin/node ${NODE_HOME}/bin/npm run build:prod"
            }
        }
        stage('Build Backend') {
            steps {
                echo "Building backend..."
                sh "cd ../crawlab"
                sh "pyenv activate crawlab"
                sh "pip install -r requirements.txt"
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
            }
        }
    }
}