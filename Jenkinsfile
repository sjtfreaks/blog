pipeline {
    agent any
    stages {
        stage ('编译镜像'){
            steps {
                echo "编译镜像"
                sh 'docker build -t docker-hugo:latest .'
                echo "结束 end"
            }
        }
        stage ('部署镜像'){
            steps {
                echo "部署镜像"
                sh 'docker stop docker-hugo'
                sh 'docker rm docker-hugo'
                sh 'docker run --name docker-hugo -d -p 4312:80 --restart=always docker-hugo:latest'
                echo "结束 end"
            }
        }
        stage ('上传制品库'){
            steps {
                echo "上传镜像到制品库"
                sh 'docker tag docker-hugo:latest hub.jobcher.com/blog/hugo:latest'
                sh 'docker push hub.jobcher.com/blog/hugo:latest'
                echo "结束 end"
            }
            post {
                success {

                    script {
                        env.COMMIT_MESSAGE = sh(script:"git --no-pager show -s -n 1 --format='%B' ${GIT_COMMIT}", returnStdout: true).trim()
                        def message = """
                            # 构建 ${BUILD_DISPLAY_NAME}
                            构建说明: ${env.COMMIT_MESSAGE}
                        """

                        dingtalk (
                            robot: '23bec93a-babe-486e-8f2f-f9486a6aac91',
                            type: 'MARKDOWN',
                            title: '流水线执行成功',
                            text: [
                                message
                            ],
                            at: [
                                '13250936269'
                            ]
                        )
                    }
            

                }
            }
        }
    }
}