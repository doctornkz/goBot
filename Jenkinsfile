node {
    agent {
        dockerfile true
    }
    stage('Build') {
        echo 'Building....'
        checkout scm
        sh 'ls -la'
        sh 'cd docker'
        docker.build("doctornkz:gobot")
    }
    stage('Test') {
        echo 'Building....'
    }
    stage('Deploy') {
        echo 'Deploying....'
    }
}