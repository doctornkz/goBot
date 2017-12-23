node {
    stage('Checkout') {
            
        echo 'Github checkout....'
        checkout scm
        sh 'pwd'
    }
    stage('GoBuild') {
        def code = pwd()
        def root = tool name: 'go1.6.2', type: 'go'
        withEnv(["GOROOT=${root}", "GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/", "PATH+GO=${root}/bin"]){
        env.PATH="${GOPATH}/bin:$PATH"
        sh 'go get -u github.com/golang/dep/...'
        sh 'dep init'
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