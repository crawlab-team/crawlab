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
                        env.TAG = 'develop'
                        env.DOCKERFILE = 'Dockerfile.local'
                    } else if (env.GIT_BRANCH == 'master') {
                        env.TAG = 'master'
                        env.DOCKERFILE = 'Dockerfile.local'
                    } 
                }
            }
        }
        stage('Build') {
            steps {
                echo "Building..."
                sh """
                docker build -t tikazyq/crawlab:${ENV:TAG} -f ${ENV:DOCKERFILE} .
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
                # 重启docker compose
                cd ./jenkins/${ENV:GIT_BRANCH}
                docker-compose down | true
                docker-compose up -d | true
                """
            }
        }
        stage('Cleanup') {
            steps {
                echo 'Cleanup...'
                sh """
                docker rmi -f `docker images | grep '<none>' | grep -v IMAGE | awk '{ print \$3 }' | xargs`
                """
            }
        }
    }
}