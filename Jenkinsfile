pipeline {
  agent { label 'x86' }

  environment {
    DOCKER_TAG = GetDockerTag()
  }

  stages {
    stage('Docker: Registry login') {
      steps {
        withCredentials([usernamePassword(credentialsId: 'bin-admin', passwordVariable: 'DOCKER_PASSWORD', usernameVariable: 'DOCKER_USERNAME')]) {
          sh 'echo "${DOCKER_PASSWORD}" | sudo docker login -u ${DOCKER_USERNAME} --password-stdin https://docker.bin.dini-mueter.net'
        }
      }
    }
    stage('Docker: Build and push branch tag') {
      steps {
        withCredentials([usernamePassword(credentialsId: 'bin-admin', passwordVariable: 'BIN_PASSWORD', usernameVariable: 'BIN_USERNAME')]) {
          ansiColor('xterm') {
            sh 'sudo docker build --build-arg VCS_REF=`git rev-parse --short HEAD` --build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` --build-arg "BIN_USER=${BIN_USERNAME}" --build-arg "BIN_PASS=${BIN_PASSWORD}" -t docker.bin.dini-mueter.net/simu/go-wiki:${DOCKER_TAG} .'
            sh 'sudo docker push docker.bin.dini-mueter.net/simu/go-wiki:${DOCKER_TAG}'
          }
        }
      }
    }
  }
}

def GetDockerTag() {
  return env.BRANCH_NAME == "master" ? "latest" : env.BRANCH_NAME.replaceAll('/','-')
}
