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
        stage('Prepare tank') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    script {
                        echo cluster_ip
                        echo tank_ip

                        def remote = [:]
                        remote.name = 'yandex-tank'
                        remote.host = tank_ip
                        remote.user = 'yc-user'
                        remote.identityFile = '/var/lib/jenkins/.ssh/id_rsa'
                        remote.allowAnyHosts = true

                        def tests_file = 'hw-' + HW_NUM + '/tests.py'

                        sshPut remote: remote, from: 'requirements.txt', into: "requirements.txt"
                        sshPut remote: remote, from: tests_file, into: "tests.py"

                        sshCommand remote: remote, sudo: true, command: 'apt install python3-pip'
                        sshCommand remote: remote, command: 'pip3 install -r requirements.txt'
                    }
                }
            }
        }
        stage('Destroy resources') {
            steps {
                sh 'bash ./destroy-k8s.sh general'
                //sh 'bash ./destroy-tank.sh general'
            }
        }
     }
}
