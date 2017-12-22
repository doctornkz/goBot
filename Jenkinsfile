node {
    stage('Build') {
        echo 'Building....'
        checkout scm
        sh 'ls -la'
        sh 'cd docker'
        sh 'ls -la'
        
        docker.build("doctornkz/gobot")
    }
    stage('Test') {
        echo 'Building....'
    }
    stage('Deploy') {
        echo 'Deploying....'
    }
}