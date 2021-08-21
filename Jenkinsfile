def cluster_ip = ""
def tank_ip = ""

pipeline {
    agent any
    options {
        skipDefaultCheckout(true)
    }
    parameters {
        string(name: 'HW_NUM', defaultValue: '0', description: 'Number of HW')
        string(name: 'TESTED_REPO', defaultValue: '', description: 'Link to github repo')
    }
    stages {
        stage('Sanity check') {
            steps {
                script {
                    if (TESTED_REPO == '') {
                        currentBuild.result = 'ABORTED'
                        error("No TESTED_REPO variable")
                    }
                    cleanWs()
                    checkout scm
                    echo "Building ${env.JOB_NAME}..."
                }
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
                        def conftest_file = 'hw-' + HW_NUM + '/conftest.py'

                        sshPut remote: remote, from: 'requirements.txt', into: "requirements.txt"
                        sshPut remote: remote, from: tests_file, into: "tests.py"
                        sshPut remote: remote, from: conftest_file, into: "conftest.py"

                        sh 'sleep 30'

                        sshCommand remote: remote, sudo: true, command: 'apt-get update'
                        sshCommand remote: remote, sudo: true, command: 'apt-get install python3-pip -y'
                        sshCommand remote: remote, sudo: true, command: 'pip3 install -r requirements.txt'
                    }
                }
            }
        }
        stage('Prepare cluster') {
            steps {
                script {
                    sh 'git clone ${TESTED_REPO} cluster-repo'
                    sh 'kubectl apply -f cluster-repo/k8s-cluster.yaml'

                    sh 'sleep 60'
                }
            }
        }
        stage('Run tests') {
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

                            def test_sh = 'pytest --junitxml=/workdir/junit.xml -q --workdir /workdir --cluster_ip ' + cluster_ip + ' tests.py'

                            sshCommand remote: remote, sudo: true, command: test_sh
                    }
                }
            }
        }
        stage('Grab test results') {
            steps {
                script {
                    def remote = [:]
                    remote.name = 'yandex-tank'
                    remote.host = tank_ip
                    remote.user = 'yc-user'
                    remote.identityFile = '/var/lib/jenkins/.ssh/id_rsa'
                    remote.allowAnyHosts = true

                    sshGet remote: remote, from: '/workdir/tank-results.json', into: "tank-results.json"
                    sshGet remote: remote, from: '/workdir/junit.xml', into: "junit.xml"

                    archiveArtifacts artifacts: 'tank-results.json', fingerprint: true
                    junit 'junit.xml'
                }
            }
        }
     }

    post {
        always {
            sh 'bash ./destroy-k8s.sh general'
            sh 'bash ./destroy-tank.sh general'
        }
    }
}
