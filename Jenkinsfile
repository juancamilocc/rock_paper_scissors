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
    - name: jnlp
      image: jenkins/inbound-agent
      resources:
        limits:
          memory: "1Gi"
          cpu: "512m"
        requests:
          memory: "500Mi"
          cpu: "256m"
  volumes:
    - name: docker-graph-storage
      emptyDir: {}
            """
        }
    }
    environment {
        REPOSITORY = 'github.com/juancamilocc/rock_paper_scissors.git' 
        BRANCH = 'deployment'
        MANIFEST = 'deployment.yaml' 
        IMAGE_TAG = sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim()
        DATE = sh(script: 'TZ="America/Bogota" date "+%Y-%m-%d-%H-%M-%S"', returnStdout: true).trim()
        LAST_CHANGE = sh(script: 'git log -1 --name-status --pretty=format:"%h %s"', returnStdout: true).trim()
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

                                withCredentials([
                                    usernamePassword(
                                        credentialsId: 'credentials-dockerhub', 
                                        usernameVariable: 'DOCKERHUB_USERNAME', 
                                        passwordVariable: 'DOCKERHUB_PASSWORD'
                                    )
                                ]) {
                                    sh '''
                                        echo $DOCKERHUB_PASSWORD | docker login -u $DOCKERHUB_USERNAME --password-stdin                      
                                        docker build -t juancamiloccc/rps-game:$IMAGE_TAG-$DATE-staging . 2> logs-docker.txt
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
                                withCredentials([
                                    usernamePassword(
                                        credentialsId: 'credentials-github', 
                                        usernameVariable: 'GIT_USERNAME', 
                                        passwordVariable: 'GIT_PASSWORD'
                                    )
                                ]) {
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
            slackSend(
                channel: 'notifications',
                color: '#00FF00',
                message: "Build of Rock Paper Scissors was successful!",
                attachments: [
                    [
                        title: "Build of Rock Paper Scissors was successful!",
                        text: "Build details",
                        fields: [
                            [title: "Date", value: "${DATE}", short: true],
                            [title: "Status", value: "Success", short: true],
                            [title: "Changes made by", value: "camilocolorado44@gmail.com", short: true],
                            [title: "Last Merge/commit", value: "${LAST_CHANGE}", short: true],
                            [title: "Project Tag", value: "rps-game:${IMAGE_TAG}-${DATE}-staging", short: true]
                        ],
                        footer: "Jenkins",
                        ts: env.BUILD_TIMESTAMP,
                        color: "#36a64f"
                    ]
                ]
            )   
        }
        failure {           
            slackSend (
                channel: 'notifications',
                color: '#00FF00',
                message: "Build of Rock Paper Scissors failed!",
                attachments: [
                    [
                        title: "Build of Rock Paper Scissors failed!",
                        text: "Build details",
                        fields: [
                            [title: "Date", value: "${DATE}", short: true],
                            [title: "Status", value: "Failure", short: true],
                            [title: "Changes made by", value: "camilocolorado44@gmail.com", short: true],
                            [title: "Last Merge/commit", value: "${LAST_CHANGE}", short: true]
                        ],
                        footer: "Jenkins",
                        ts: env.BUILD_TIMESTAMP,
                        color: "#ff0000"
                    ]
                ]
            )
        }
    }
}
