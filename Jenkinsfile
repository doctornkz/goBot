node {

    stage('Build'){

    def root = tool name: 'go1.8', type: 'go'
    ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/github.com/doctornkz/goBot") {
        
        withEnv(["GOROOT=${root}", "GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/", "PATH+GO=${root}/bin"]) {
            env.PATH="${GOPATH}/bin:$PATH"
        
            git url: 'https://github.com/doctornkz/goBot.git'
            sh 'go version'
            sh 'go get -u github.com/golang/dep/...'
            sh 'dep init'
            sh 'go build main.go'

            }
            //stage 'Test'
            //sh 'go vet'
            //sh 'go test -cover'
            
           //sh 'go build .'
            
            stage ('Deploy'){
            // Do nothing.
            }
        }
    }
}