def dockerfileDict = [

    user-check_api: [
        filename: "group_check.Dockerfile"
    ]
]


def dockerfileDictList = dockerfileDict.keySet().toList()

pipeline {

    agent { label "default-agent-us" }

    options {
        durabilityHint('PERFORMANCE_OPTIMIZED')
        ansiColor('xterm')
        timestamps()
        buildDiscarder logRotator(artifactDaysToKeepStr: '14', artifactNumToKeepStr: '', daysToKeepStr: '14', numToKeepStr: '')
        skipStagesAfterUnstable()
        parallelsAlwaysFailFast()
        disableConcurrentBuilds()
    }

    parameters { choice(name: 'DOCKERFILE_DATA', choices: dockerfileDictList, description: '') }

    environment {
        DOCKERFILE_NAME = "${dockerfileDict[params.DOCKERFILE_DATA]['filename']}"
    }

    stages {
        stage('Push to docker') {
            steps {
                script {
                    dockerImageName = pushDockerfile (
                            pushDockerfileBuildDir: env.DOCKERFILE_DIR,
                            pushDockerfileName: env.DOCKERFILE_NAME,
                            pushDockerfileCredentialsId: "credentials",
                            pushDockerfileOnlyBuild: false)
                }
            }
        }
    }
    post {

        // cleanup
        cleanup {
            deleteDir()
        }
    }

}