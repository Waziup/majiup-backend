pipeline {
    
    agent any

    options {
        timeout(time: 1, unit: 'HOURS')
    }

    stages {
        stage('Checkout') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    // Clone the frontend repository
                    sh 'git clone https://github.com/Waziup/majiup-waziapp.git'
                }
            }
        }

        stage ('Build') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    // Navigate to the cloned frontend repository and build
                    dir('majiup-waziapp') {
                        // Build the majiup-frontend
                        sh 'npm install --legacy-peer-deps'
                        sh 'npm run build'
                        sh 'cp -r dist/ serve/'
                    }
                    sh 'sudo docker build --platform linux/arm64  -t waziup/majiup .'
                                
                }
            }
        }

        stage('Deploy') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {
                    sh 'docker push waziup/majiup'
                    sh 'remote_start_waziapp.sh'
                    
                }
            }
        }

        stage('Run Tests') {
            steps {
                catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {

                    // Navigate to the tests folder and run tests
                    dir('tests') {
                        sh 'python tests.py'
                    }
                }
            }
        }        
    }

    post {
        always {
            junit 'tests/test_results.xml'
        }
    }
}
