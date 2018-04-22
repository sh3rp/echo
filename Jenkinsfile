node {   
    ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/") {
        def root = tool name: 'Go 1.10', type: 'go'

        withEnv(["GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"]) {
            env.PATH="${GOPATH}/bin:${root}/bin:$PATH"
            env.TIME = sh(returnStdout: true, script: "date").trim()
            env.GOVERSION = sh(returnStdout: true, script: "go version | awk '{print \$3}'").trim()

            stage('Checkout'){
                echo 'Checking out SCM'
                sh 'go get -u github.com/sh3rp/echo'
                env.COMMIT = sh(returnStdout: true, script: "cd ${GOPATH}/src/github.com/sh3rp/echo && git rev-list -1 HEAD").trim()
            }
            
            stage('Pre Build'){
                echo 'Pulling Dependencies'
        
                sh 'go version'
                sh 'go get -u github.com/golang/dep/cmd/dep'                
                sh 'cd ${GOPATH}/src/github.com/sh3rp/echo && dep ensure' 
            }
        
            stage('Build'){
                echo 'Building Executable'
                echo "Time: ${TIME}"
                echo "Commit: ${COMMIT}"
                echo "Version: ${GOVERSION}"
            
                sh """\
                    cd $GOPATH/src/github.com/sh3rp/echo && \
                    go build -ldflags -s -X main.BuildTime="$TIME" -X main.GitCommit="$COMMIT" -X main.GoVersion="$GOVERSION" \
                """
            }
            
        }
    }
}