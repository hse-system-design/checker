pipeline {
    agent any
    parameters {
        string(name: 'HW_NUM', defaultValue: '0', description: 'Number of HW')
    }
    stages {
        stage('Hello') {
            steps {
                echo 'Hello World'
            }
        }
        stage('Run sh') {
            steps {
                sh 'bash ./run.sh ${HW_NUM}'
            }
        }
    }
}
