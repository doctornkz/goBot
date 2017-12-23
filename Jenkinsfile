node {
    stage('Checkout') {
            
        echo 'Github checkout....'
        checkout scm
        sh 'pwd'
    }
    stage('GoBuild') {
        def code = pwd()
        def root = tool name: 'go1.6.2', type: 'go'
        withEnv(["GOROOT=${root}", "GOPATH=${root}","PATH+GO=${root}"]){
        
        sh "ls -la ${code}"
        sh "ls -la ${root}"
        
        sh 'go get github.com/tools/godep'

        
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