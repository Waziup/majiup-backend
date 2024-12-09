pipeline {
    
    agent any

    options {
        timeout(time: 1, unit: 'HOURS')
    }

    stages {
 
        stage ('Build') {
            steps {
                sh 'sudo docker buildx build --tag=waziup/majiup:latest . --platform=linux/arm64 --load --no-cache --progress plain'
            }
        }

        stage('Deploy') {
            steps {
               sh 'docker push waziup/majiup:latest'
               sh 'sudo chmod +x ./remote_start_waziapp.sh'
               sh 'sudo ./remote_start_waziapp.sh'
            }
        }

        // stage('Run Tests') {
        //     steps {
        //         catchError(buildResult: 'SUCCESS', stageResult: 'FAILURE') {

        //             // Navigate to the tests folder and run tests
        //             dir('tests') {
        //                 sh 'python tests.py'
        //             }
        //         }
        //     }
        // }        
    }

    // post {
    //     always {
    //         junit 'tests/test_results.xml'
    //     }
    // }
}
