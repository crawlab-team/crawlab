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
                        env.MODE = 'develop'
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
                docker build -t tikazyq/crawlab:latest -f Dockerfile.local .
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
                cd ./jenkins
                docker-compose stop | true
                docker-compose up -d
                """
            }
        }
    }
}