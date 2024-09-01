pipeline {
    agent {
        kubernetes {
        cloud 'kubernetes-staging'    
        defaultContainer 'jnlp'
        yaml """
apiVersion: v1
kind: Pod
metadata:
  name: rocky-pod
  namespace: jenkins
spec:
  containers:
    - name: rocky
      image: ghcr.io/juancamilocc/builders:rocky8-docker
      imagePullPolicy: IfNotPresent
      tty: true
      securityContext:
        runAsUser: 0
        privileged: true
      resources:
        limits:
          memory: "2Gi"
          cpu: "750m"
        requests:
          memory: "1Gi"
          cpu: "500m"
      volumeMounts:
        - name: docker-graph-storage
          mountPath: /var/lib/docker
  volumes:
    - name: docker-graph-storage
      emptyDir: {}
            """
            containerTemplate {
            name 'jnlp'
            image 'jenkins/inbound-agent'
            resourceRequestCpu '256m'
            resourceRequestMemory '500Mi'
            resourceLimitCpu '512m'
            resourceLimitMemory '1000Mi'
            }
        }
    }
    environment {
        REPOSITORY = 'github.com/juancamilocc/rock_paper_scissors.git' 
        BRANCH = 'deployment'
        MANIFEST = 'deployment.yaml' 
        IMAGE_TAG = sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim()
        DATE = sh(script: 'TZ="America/Bogota" date "+%Y-%m-%d-%H-%M-%S"', returnStdout: true).trim()
        RETRY_COUNTS = 2
    }
    stages {
        stage('Build and Push image') {
            steps {
                container('rocky') {
                    script {
                        retry(RETRY_COUNTS) {
                            try {
                                sh 'git config --global --add safe.directory $WORKSPACE'

                                withCredentials([usernamePassword(credentialsId: 'credentials-dockerhub', usernameVariable: 'DOCKERHUB_USERNAME', passwordVariable: 'DOCKERHUB_PASSWORD')]) {
                                    sh '''
                                        echo $DOCKERHUB_PASSWORD | docker login -u $DOCKERHUB_USERNAME --password-stdin                      
                                        docker build -t juancamiloccc/rps-game:$IMAGE_TAG-$DATE-staging .
                                        docker push juancamiloccc/rps-game:$IMAGE_TAG-$DATE-staging
                                    '''
                                }

                            } catch (Exception e) {
                                echo "Error occurred: ${e.message}"
                                echo "Retrying..."
                                error("The stage 'Build and Push image' failed")
                            }
                        }
                    }
                }
            }
        }
        stage('Update deployment') {
            steps {
                container('rocky') {
                    script {
                        retry(RETRY_COUNTS) {
                            try {
                                withCredentials([usernamePassword(credentialsId: 'credentials-github', usernameVariable: 'GIT_USERNAME', passwordVariable: 'GIT_PASSWORD')]) {
                                    sh '''
                                        git config --global user.email "jccoloradoc@uqvirtual.edu.co"
                                        git config --global user.name "juancamilocc"
                                        git clone -b $BRANCH --depth 5 https://GIT_USERNAME:$GIT_PASSWORD@$REPOSITORY
                                        cd rock_paper_scissors/deployment/staging
                                        sed -i "s/\\(image:.*:\\).*/\\1$IMAGE_TAG-$DATE-staging/" $MANIFEST
                                        git add $MANIFEST 
                                        git commit -m "Trigger Build"                     
                                        git push origin $BRANCH
                                    '''

                                    // Delete repository
                                    sh 'rm -rf rock_paper_scissors'
                                }
                            } catch (Exception e) {
                                echo "Error occurred: ${e.message}"
                                echo "Retrying..."
                                error("The stage 'Update Deployment' failed")
                            }
                        }    
                    }
                }
            }
        }
    }
    post {
        success {
            slackSend channel: 'jenkins-notifications',
                message: "Build ${currentBuild.currentResult} for pipeline ${env.JOB_NAME} build #${env.BUILD_NUMBER} ${BUILD_URL}"   
        }
        failure {
            echo "FAILURE"
        }
    }
}
