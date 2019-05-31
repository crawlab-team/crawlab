pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo "Git branch: ${env.GIT_BRANCH}"
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