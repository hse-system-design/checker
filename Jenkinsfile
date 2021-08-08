def CLUSTER_IP = ""

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
        stage('Deploy k8s cluster') {
            steps {
                CLUSTER_IP = sh(script: 'bash ./deploy-k8s.sh general ./hw-${HW_NUM}/cluster-config.json', returnStdout: true)
                echo CLUSTER_IP
            }
        }
    }
}
