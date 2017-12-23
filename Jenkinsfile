node {

    stage('Go Build'){

    def root = tool name: 'go1.8', type: 'go'
    ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/github.com/doctornkz/goBot/") {
        
        withEnv(["GOROOT=${root}", "GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/", "PATH+GO=${root}/bin"]) {
            env.PATH="${GOPATH}/bin:$PATH"
        
            git url: 'https://github.com/doctornkz/goBot.git'
            sh 'go version'
            sh 'go get -u github.com/golang/dep/...'
            sh 'dep init'
            sh 'go build -o goBot '
        }}
    stage ('Docker build'){
            ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/github.com/doctornkz/goBot/") {
            sh "cp goBot docker"
            sh 'cd docker'
            sh 'docker build -t doctornkz/gobot docker'
            sh 'docker push doctornkz/gobot'

            }

            // Do nothing.
            
            
            }
        
    }
}