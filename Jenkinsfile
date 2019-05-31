pipeline {
    agent {crawlab}

    stages {
        stage('Build Frontend') {
            steps {
                echo "Building frontend..."
                sh "cd frontend"
                sh "npm install -g yarn pm2"
                sh "yarn install"
                sh "npm run build:prod"
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