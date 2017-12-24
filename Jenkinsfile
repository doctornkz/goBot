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
            }
    }
        stage('Push image') {
        /* Finally, we'll push the image with two tags:
         * First, the incremental build number from Jenkins
         * Second, the 'latest' tag.
         * Pushing multiple tags is cheap, as all the layers are reused. */
            docker.withRegistry('https://registry.hub.docker.com', 'docker-hub') {
            docker.push("doctornkz/gobot")
            
        }
//            sh 'docker build -t doctornkz/gobot docker'
//            sh 'docker push doctornkz/gobot'
            
            }
        
    }
}