node {
    stage('Checkout') {
        echo 'Github checkout....'
        checkout scm
    }
    stage('GoBuild') {
        env.GOPATH = './'
        echo 'Go building....'
        sh 'go build main.go'
    }

    stage('DockerBuild'){
        dir('docker'){
            docker.build("doctornkz/gobot")
    
        }

    }   
    
    stage('Test') {
        echo 'Building....'
    }

    stage('Deploy') {
        echo 'Deploying....'
    }
}