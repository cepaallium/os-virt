def modules = """
[
    {
        "osServiceName": "os-virt",
        "path": "${SERVICE_GROUP}/os-virt",
        "version": "${BUILD_VERSION}",
        "name": "os-virt-${BUILD_VERSION}.tar",
        "localPath": "./helm-install/${SERVICE_NAME}-${BUILD_VERSION}/images/",
        "pipelineId": "1753246984552054784",
        "parameters": {"TRIGGER_GIT_REF": "${TRIGGER_GIT_REF}","BUILD_VERSION": "${BUILD_VERSION}","OS_ARCH": "${OS_ARCH}"}
    }
]
"""

sh """
mkdir -p ./helm-install/${SERVICE_NAME}-${BUILD_VERSION}/images/

cd ./helm-install/${SERVICE_NAME}-${BUILD_VERSION}/images/

"""
pullArtifactBatchWithGitTrigger modules: modules, timeout: 150