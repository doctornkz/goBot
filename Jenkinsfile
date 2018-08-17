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
            sh 'go build -ldflags \\"-X main.currentVersion=currentMaster\\" -o goBot '
            sh 'cp goBot docker'
        }}}

    stage('Docker build'){
            ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src/github.com/doctornkz/goBot/docker") {
            sh 'docker build -t doctornkz/gobot .'

            def gobotImage = docker.build("doctornkz/gobot")
                docker.withRegistry('https://registry.hub.docker.com', 'docker-hub') {
                gobotImage.push("latest")
        }
    }
}
}