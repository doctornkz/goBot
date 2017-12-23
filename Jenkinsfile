node {
    stage('Checkout') {
            
        echo 'Github checkout....'
        checkout scm
        sh 'pwd'
    }
    stage('GoBuild') {
        def code = pwd()
        def root = tool name: 'go1.6.2', type: 'go'
        withEnv(["GOROOT=${root}", "GOPATH=${root}/bin","PATH+GO=${root}/bin"]){
        
        sh "ls -la ${code}"
        sh "ls -la ${root}"
        sh 'go get github.com/tools/godep'
        sh "mv ${root}/bin/bin/godep ${root}/bin/"
        sh 'pwd'
        sh 'godep save ./...'
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