node {   
    ws("${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/") {
        def root = tool name: 'Go 1.10', type: 'go'

        withEnv(["GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}"]) {
            env.PATH="${GOPATH}/bin:${root}/bin:$PATH"
            
            stage('Checkout'){
                echo 'Checking out SCM'
                sh 'go get -u github.com/sh3rp/echo'
            }
            
            stage('Pre Test'){
                echo 'Pulling Dependencies'
        
                sh 'go version'
                sh 'go get -u github.com/golang/dep/cmd/dep'
                sh 'go get -u github.com/golang/lint/golint'
                
                sh 'cd ${GOPATH}/src/github.com/sh3rp/echo && dep ensure' 
            }
    
            stage('Test'){
                sh 'cd $GOPATH && go list ./... | grep -v /vendor/ | grep -v github.com | grep -v golang.org > projectPaths'
                
                def paths = sh returnStdout: true, script: """awk '\$0="./src/"\$0' projectPaths"""
                
                echo 'Testing'
                sh """cd $GOPATH && go test -race -cover ${paths}"""
            }
        
            stage('Build'){
                echo 'Building Executable'
            
                sh """cd $GOPATH/src/sh3rp/echo && go build -ldflags '-s'"""
            }
            
        }
    }
}