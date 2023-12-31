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
                // Navigate to the cloned frontend repository and build
                //dir('majiup-waziapp') {
                //    // Build the majiup-frontend
                //    sh 'pnpm install'
                //    sh 'pnpm build'
                //    sh 'cp -r dist/ serve/'
                //}
                sh 'sudo docker buildx bake --load --progress plain'
            }
        }

        stage('Deploy') {
            steps {
               sh 'docker push waziup/majiup'
               sh 'remote_start_waziapp.sh'
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
