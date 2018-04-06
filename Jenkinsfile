pipeline {
    agent { docker 'iron/go:dev' }

    environment { HOME='.' }

    stages {
        stage('Prepare') {
            steps {
                echo 'Installing dependencies..'
                sh 'go get -u github.com/gorilla/mux'
                sh 'go get -u github.com/go-sql-driver/mysql'
            }
        }
        stage('Test') {
            steps {
                echo 'Running tests'
                sh 'go test ./...'
            }
        }
        stage('Build') {
            steps {
                echo 'Building sources'
		sh 'go build trader.go trader.config.go'
            }
        }
    }
}
