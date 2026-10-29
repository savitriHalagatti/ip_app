// CI/CD pipeline for the FastAPI app (ip_app/fast_api)
//
// Stages: Checkout -> Install -> Lint -> Unit Test -> Security Scan -> Build Image
//         -> Container Test -> Push Image (optional) -> Deploy (dev)
//
// Jenkins job: Pipeline script from SCM, Script Path = fast_api/Jenkinsfile
// Jenkins agent needs: python3 + python3-venv, docker (jenkins user in the docker group), curl
// Credentials (only for Push Image): dockerhub-creds = Docker Hub username + access token

pipeline {
    agent any

    parameters {
        booleanParam(name: 'PUSH_IMAGE', defaultValue: false, description: 'Push the image to Docker Hub (needs dockerhub-creds)')
        booleanParam(name: 'DEPLOY', defaultValue: true, description: 'Run the new image on this Jenkins agent (dev)')
    }

    environment {
        APP_DIR     = 'fast_api'
        IMAGE       = 'fastapi'                 // same name as in docker-compose.yml
        TAG         = "${env.BUILD_NUMBER}"
        DEPLOY_PORT = '9080'                    // dev app: http://<jenkins-agent-ip>:9080/
        TEST_HOST   = 'localhost'               // if Jenkins itself runs in Docker: host.docker.internal
    }

    options {
        buildDiscarder(logRotator(numToKeepStr: '10'))
        disableConcurrentBuilds()
        timeout(time: 30, unit: 'MINUTES')
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Install') {
            steps {
                dir(env.APP_DIR) {
                    sh '''
                        python3 -m venv .venv
                        . .venv/bin/activate
                        pip install -q --upgrade pip
                        pip install -q -r requirements-dev.txt
                    '''
                }
            }
        }

        stage('Lint') {
            steps {
                dir(env.APP_DIR) {
                    // ruff check = code problems, ruff format --check = code style
                    sh '. .venv/bin/activate && ruff check . && ruff format --check .'
                }
            }
        }

        stage('Unit Test') {
            steps {
                dir(env.APP_DIR) {
                    // fails if coverage is below 80%
                    sh '''
                        . .venv/bin/activate
                        pytest --junitxml=reports/unit-tests.xml \
                               --cov=app --cov-report=term-missing \
                               --cov-report=xml:reports/coverage.xml --cov-report=html:reports/coverage-html \
                               --cov-fail-under=80
                    '''
                }
            }
            post {
                always {
                    junit "${env.APP_DIR}/reports/unit-tests.xml"
                    archiveArtifacts artifacts: "${env.APP_DIR}/reports/coverage-html/**", allowEmptyArchive: true
                }
            }
        }

        stage('Security Scan') {
            steps {
                dir(env.APP_DIR) {
                    // fails if a library in requirements.txt has a known vulnerability
                    sh '. .venv/bin/activate && pip-audit -r requirements.txt'
                }
            }
        }

        stage('Build Image') {
            steps {
                dir(env.APP_DIR) {
                    sh 'docker build -t $IMAGE:$TAG -t $IMAGE:latest .'
                }
            }
        }

        stage('Container Test') {
            steps {
                dir(env.APP_DIR) {
                    // run the new image on a free port and test it over HTTP
                    sh '''
                        docker rm -f $IMAGE-test 2>/dev/null || true
                        docker run -d --name $IMAGE-test -p 8000 $IMAGE:$TAG
                        PORT=$(docker port $IMAGE-test 8000/tcp | head -1 | sed 's/.*://')
                        for i in $(seq 1 30); do curl -sf http://$TEST_HOST:$PORT/greet > /dev/null && break; sleep 1; done
                        . .venv/bin/activate
                        BASE_URL=http://$TEST_HOST:$PORT pytest tests/integration --junitxml=reports/integration-tests.xml
                    '''
                }
            }
            post {
                always {
                    junit "${env.APP_DIR}/reports/integration-tests.xml"
                    sh 'docker rm -f $IMAGE-test || true'
                }
            }
        }

        stage('Push Image') {
            when { expression { params.PUSH_IMAGE } }
            steps {
                withCredentials([usernamePassword(credentialsId: 'dockerhub-creds',
                        usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                    sh '''
                        echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
                        docker tag $IMAGE:$TAG $DOCKER_USER/$IMAGE:$TAG
                        docker tag $IMAGE:$TAG $DOCKER_USER/$IMAGE:latest
                        docker push $DOCKER_USER/$IMAGE:$TAG
                        docker push $DOCKER_USER/$IMAGE:latest
                        docker logout
                    '''
                }
            }
        }

        stage('Deploy (dev)') {
            when { expression { params.DEPLOY } }
            steps {
                sh '''
                    docker rm -f $IMAGE-dev 2>/dev/null || true
                    docker run -d --name $IMAGE-dev --restart unless-stopped -p $DEPLOY_PORT:8000 $IMAGE:$TAG
                    for i in $(seq 1 30); do curl -sf http://$TEST_HOST:$DEPLOY_PORT/greet > /dev/null && break; sleep 1; done
                    curl -sf http://$TEST_HOST:$DEPLOY_PORT/greet
                '''
            }
        }
    }

    post {
        success {
            echo "Image ${env.IMAGE}:${env.TAG} passed all checks. Dev app: http://<jenkins-agent-ip>:${env.DEPLOY_PORT}/"
        }
        always {
            sh 'docker image prune -f || true'
        }
    }
}
