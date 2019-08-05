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
                echo ${ENV:GIT_BRANCH}
                """
                sh """
                cd ./jenkins/${ENV:GIT_BRANCH}
                docker-compose stop | true
                docker-compose up -d --scale worker=3
                """
            }
        }
        stage('Cleanup') {
            steps {
                echo 'Cleanup...'
                sh """
                docker rmi `docker images | grep '<none>' | grep -v IMAGE | awk '{ print \$3 }' | xargs` | true
                docker rm `docker ps -a | grep Exited | awk '{ print \$1 }' | xargs` | true
                """
            }
        }
    }
}