node {
    stage('Checkout') {
        echo 'Github checkout....'
        checkout scm
    }
    stage('GoBuild') {
        def root = tool name: 'Go 1.6.2', type: 'go'
        withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) 
        sh 'go version'
    
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