node {   
    ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/") {
        def root = tool name: 'Go 1.10', type: 'go'

        withEnv(["GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"]) {
            env.PATH="${GOPATH}/bin:${root}/bin:$PATH"
            
            stage('Checkout'){
                echo 'Checking out SCM'
                sh 'go get -u github.com/sh3rp/echo'
            }
            
            stage('Pre Build'){
                echo 'Pulling Dependencies'
        
                sh 'go version'
                sh 'go get -u github.com/golang/dep/cmd/dep'                
                sh 'cd ${GOPATH}/src/github.com/sh3rp/echo && dep ensure' 
            }
        
            stage('Build'){
                echo 'Building Executable'
            
                sh """\
                    export TIME=$(date) && \
                    export COMMIT=$(git rev-list -1 HEAD) && \
                    export GOVER=$(go version) && \
                    cd $GOPATH/src/github.com/sh3rp/echo && \
                    go build -ldflags '-s -X main.BuildTime=${TIME} -X main.GitCommit=${COMMIT} -X main.GoVersion=${GOVER}' \
                """
            }
            
        }
    }
}