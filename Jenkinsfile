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
                        env.TAG = 'develop'
                        env.BASE_URL = '/dev'
                    } else if (env.GIT_BRANCH == 'master') {
                        env.MODE = 'production'
                        env.TAG = 'master'
                        env.BASE_URL = '/demo'
                    } 
                }
            }
        }
        stage('Build') {
            steps {
                echo "Building..."
                sh """
                docker build -t tikazyq/crawlab:${ENV:TAG} -f Dockerfile.local .
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
                docker-compose up -d
                """
            }
        }
        stage('Cleanup') {
            steps {
                echo 'Cleanup...'
                sh """
                # remove unused containers
                container_ids=`docker ps -a | grep Exited | awk '{ print \$1 }' | xargs`
                if [ \\$container_ids -eq "" ];
                then
                    :
                else
                    docker rm \$container_ids
                fi

                # remove unused images
                image_ids=`docker images | grep '<none>' | grep -v IMAGE | awk '{ print \$3 }' | xargs`
                if [ \\$image_ids -eq "" ];
                then
                    :
                else
                    docker rmi \$image_ids
                fi
                """
            }
        }
    }
}