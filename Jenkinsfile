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
                        env.DOCKERFILE = 'docker/Dockerfile.frontend.alpine'
                    } else if (env.GIT_BRANCH == 'backend-master') {
                        env.TAG = 'master-alpine'
                        env.DOCKERFILE = 'docker/Dockerfile.master.alpine'
                    } else if (env.GIT_BRANCH == 'backend-worker') {
                        env.TAG = 'worker-alpine'
                        env.DOCKERFILE = 'docker/Dockerfile.worker.alpine'
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
                if [[ ${ENV:GIT_BRANCH} = master ]] || [[ ${ENV:GIT_BRANCH} = develop ]]; then
                    # 重启docker compose
                    cd ./jenkins/${ENV:GIT_BRANCH}
                    docker-compose stop | true
                    docker-compose up -d
                else
                    # 发布到Dockerhub
                    docker push tikazyq/crawlab:${ENV:TAG}
                fi
                """
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