node {
    stage('Checkout') {
        echo 'Github checkout....'
        checkout scm
    }
    stage('GoBuild') {
        def root = tool name: 'go1.6.2', type: 'go'
        withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]){
        sh 'go install github.com/tools/godep'
        sh 'godep get  ./...'
        sh 'go build main.go'

        }
    }

    stage('DockerBuild'){
        sh 'cp goBot docker'
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