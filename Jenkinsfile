pipeline {
    agent {
        node {
            label 'crawlab'
        }
    }

    stages {
        stage('Setup') {
            steps {
                echo "Running Setup..."
                script {
                    if (env.GIT_BRANCH == 'develop') {
                        env.MODE = 'test'
                    } else if (env.GIT_BRANCH == 'master') {
                        env.MODE = 'production'
                    } else {
                        env.MODE = 'test'
                    }
                }
            }
        }
        stage('Build') {
            steps {
                echo "Building..."
                sh """
                echo `whoami`
                docker build -t crawlab:latest .
                """
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
                sh """
                docker stop crawlab | true
                docker run -d --rm --restart always --name crawlab \
                    -p 8080:8080 \
                    -p 8000:8000 \
                    -v /home/yeqing/.env.production:/opt/crawlab/frontend/.env.production \
                    -v /home/yeqing/config.py:/opt/crawlab/crawlab/config/config.py
                """
            }
        }
    }
}