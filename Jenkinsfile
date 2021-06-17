@Library(['shared-libs']) _
 
pipeline {
    agent { docker { image 'golang' } }
 
    options {
        ansiColor('xterm')
        timestamps()
        timeout(time: 1, unit: 'HOURS')
        gitLabConnection('GitLab Master')
        buildDiscarder(logRotator(numToKeepStr: '100', artifactNumToKeepStr: '10'))
    }
 
    environment {
        HOME="${WORKSPACE}"
        PYTHONUNBUFFERED=1
    }
 
    parameters {
        string(name: 'REF', defaultValue: '\${gitlabBranch}', description: 'Commit to build')
    }
 
    stages {
        stage('Prep') {
            steps {
                script {
                    updateGitlabCommitStatus(name: 'Jenkins CI', state: 'running')
                }
            }
        }
        stage('Compile') {
            steps {
                echo "Runs in ${WORKSPACE}"
                sh "pwd"
 
                make binary
            }
        }
	stage('Test') {
            steps {
                echo "Running tests"
		// Tests require supported GPU
                // make test-main
                make check-format
            }
        }
    }
    post {
        always {
            script{
                String status = (currentBuild.currentResult == "SUCCESS") ? "success" : "failed"
                updateGitlabCommitStatus(name: 'Jenkins CI', state: status)
            }
        }
        cleanup {
            cleanWs()
        }
    }
}
