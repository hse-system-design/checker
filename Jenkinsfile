def cluster_ip = ""
def tank_ip = ""

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
                script {
                    sh 'bash ./deploy-k8s.sh general ./hw-${HW_NUM}/cluster-config.json'
                    cluster_ip = readFile(file: 'cluster_ip.txt').trim()
                    echo cluster_ip
                }
            }
        }
        stage('Deploy yandex tank') {
            steps {
                script {
                    sh 'bash ./deploy-tank.sh general /var/lib/jenkins/.ssh/id_rsa.pub'
                    tank_ip = readFile(file: 'tank_ip.txt').trim()
                    echo tank_ip
                }
            }
        }
        stage('Test') {
            steps {
                echo cluster_ip
                echo tank_ip
            }
        }
        stage('Destroy resources') {
            steps {
                sh 'bash ./destroy-k8s.sh general'
                sh 'bash ./destroy-tank.sh general'
            }
        }
     }
}
