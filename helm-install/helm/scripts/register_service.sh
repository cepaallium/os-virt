#!/bin/sh

KONG_PROXY_SVC="os-kong-proxy-svc.default.svc"
KONG_PROXY_PORT="11000"
KONG_SVC="os-kong-admin-svc.default.svc.cloudos"
KONG_PORT="11001"
SERVICE_NAME="os-virt"
SVCMGT_SVC="os-cill-svc.default.svc.cloudos"
SVCMGT_PORT="11021"
NAMESPACE="cloudos-cvs"

## 1. 检查服务是否就绪
function check_service() {
    checkResult=`curl http://${KONG_SVC}:${KONG_PORT}/read_health`
    echo "checkResult: $checkResult"
    count=0
    while [ -z "$checkResult" ]
    do
        sleep 10
        checkResult=`curl http://${KONG_SVC}:${KONG_PORT}/read_health`
        echo "checkResult: $checkResult"
        let "count+=1"
        if [ $count -gt 30 ];then
            echo "Service registered failed, please check kong service!"
            exit 1
        fi
    done
}

## 2. 获取注册请求体
function obtain_register_body() {
    if [ ! -f "/opt/shellScripts/register_service_body.json" ];then
        echo "There is no register body file,register service failed!"
        exit 10
    fi
    jq -r . /opt/shellScripts/register_service_body.json
    if [ $? -ne "0" ];then
        echo "There is something wrong with the file of register_service_body.json,please check the format of the file!"
        exit 15
    fi
    register_body=`cat /opt/shellScripts/register_service_body.json`
}


## 3. 注册服务
function registered_service() {
    serviceId=`curl -X POST \
        http://${KONG_SVC}:${KONG_PORT}/registered \
        -H 'Content-Type: application/json' \
        -d @/opt/shellScripts/register_service_body.json | jq .service.id | sed 's/"//g'`
    echo "serviceId: $serviceId"
    if [ -n serviceId ]; then
        echo "Service register successfully, please check in http://${KONG_SVC}:${KONG_PORT}/services/${serviceId}"
        register_statement
    else
        echo "Register service failed!"
        exit 20
    fi
}

function register_statement() {
    RESOURCE_KEY=`cat /opt/shellScripts/register_service_body.json | jq .service.name | sed 's/"//g'`
    for (( times=0; times<= 5; times++ ))
    do

        virtAdminToken=$(curl -s -i -X POST \
            http://${KONG_PROXY_SVC}:${KONG_PROXY_PORT}/sys/auth/v1/tokens \
            -H 'Content-Type: application/json' \
            -d '{"auth":{"identity":{"methods":["password"],"password":{"user":{"name":"virtAdmin"},"password":"Passw0rd@_"}}}}' \
            | grep X-Subject-Token | sed 's/:/\n/g' | sed '1d')

        echo "virtAdminToken: ${virtAdminToken}"

        return_code=$(curl -o /dev/null -s -w %{http_code} -X GET \
           -H 'X-Auth-Token: '"${virtAdminToken}"'' \
           http://${SVCMGT_SVC}:${SVCMGT_PORT}/v1/service/statement?namespace=${NAMESPACE}\&servicename=${SERVICE_NAME}\&key=${RESOURCE_KEY}\&type=service)

        if [ ${times} -eq 5 ];then
            echo "Statement Register Failed,${SVCMGT_SVC} Interface is not ready!"
        fi
        if [[ ${return_code} =~ ^[2-3] ]];then
            echo "Statement Register successfully!"
            break
        else
            echo "Statement Register failed!"
            sleep 5
            continue
        fi
    done
}

main() {
  check_service
  obtain_register_body
  registered_service
}
main
