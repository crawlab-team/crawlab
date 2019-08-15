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
                    } else if (env.GIT_BRANCH == 'frontend') {
                        env.TAG = 'frontend-alpine'
                        env.DOCKERFILE = 'frontend/Dockerfile.frontend.alpine'
                    } else if (env.GIT_BRANCH == 'backend-master') {
                        env.TAG = 'master-alpine'
                        env.DOCKERFILE = 'frontend/Dockerfile.master.alpine'
                    } else if (env.GIT_BRANCH == 'backend-worker') {
                        env.TAG = 'worker-alpine'
                        env.DOCKERFILE = 'frontend/Dockerfile.worker.alpine'
                    } 
                }
            }
        }
        stage('Build') {
            steps {
                echo "Building..."
                sh """
                docker build -t tikazyq/crawlab:${ENV:TAG} -f ${ENV.DOCKERFILE} .
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
                script {
                    if (env.GIT_BRANCH == 'master' || env.GIT_BRANCH == 'develop') {
                        sh """
                        cd ./jenkins/${ENV:GIT_BRANCH}
                        docker-compose stop | true
                        docker-compose up -d
                        """
                    } else {
                        sh """
                        docker push tikazyq/crawlab:${ENV:TAG}
                        """
                    }
                }
            }
        }
        stage('Cleanup') {
            steps {
                echo 'Cleanup...'
                sh """
                docker image prune -f
                """
            }
        }
    }
}