pipeline {
    agent any
    options {
        timeout(time: 1, unit: 'HOURS')
    }
    stages {
        stage('Checkout') {
            steps {
                // Clone the frontend repository
                sh 'git clone https://github.com/Waziup/majiup-waziapp.git'

                // Clone the wazigate-edge repository
                sh 'git clone https://github.com/Waziup/wazigate-edge.git'                                
            }
        }

        stage ('Build') {
            steps {
                // Navigate to the cloned frontend repository and build
                dir('majiup-waziapp') {
                    // Build the majiup-frontend
                    sh 'npm install'
                    sh 'npm run build'
                }                
            }
        }

        stage('Setup Frontend Production Files') {
            steps {
                // Copy the dist folder to the main repository folder
                sh 'cp -r majiup-waziapp/dist ../serve/'
            }
        }

        stage('Run Majiup') {
            steps {
                // Navigate to the cloned wazigate-edge
                dir('wazigate-edge') {
                    // Build the Docker image
                    sh 'docker build --tag=wazigate-edge .'
                    
                    // Run the Docker container (Waziedge)
                    sh 'docker run -d --name wazigate-edge wazigate-edge'
                }
                // Navigate to the root directory
                dir('../') {
                    // Build the Docker Compose in the main directory
                    sh 'docker-compose up -d'
                }
            }
        }
    }

    post {
        always {

        }
    }
}
