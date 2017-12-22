#!/usr/bin/env groovy
stage('Checkout') {
    // some block
checkout([$class: 'GitSCM', branches: [[name: '*/master']], doGenerateSubmoduleConfigurations: false, extensions: [], submoduleCfg: [], userRemoteConfigs: [[credentialsId: 'github.com', url: 'https://github.com/doctornkz/goBot']]])
}


