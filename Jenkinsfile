node {
    stage('Checkout') {
        def coderoot = pwd()        
        echo 'Github checkout....'
        checkout scm
    }
    stage('GoBuild') {
        def root = tool name: 'go1.6.2', type: 'go'
        withEnv(["PATH+GO=${root}", "GOPATH=${root}"]){
        sh 'go get github.com/tools/godep'
        sh "${root}/bin/godep get ./..."
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