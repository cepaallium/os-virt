#!/bin/bash

KONG_SVC="os-kong-proxy-svc.default.svc"
KONG_PORT="11000"

## 1. 检查服务是否就绪
function check_service() {
    for (( times=0; times<= 30; times++ ))
    do
        return_code=$(curl -o /dev/null -s -w %{http_code} http://${SVCMGT_SVC}:${SVCMGT_PORT}/v1/rest_health)
        if [ ${times} -eq 30 ];then
            echo "Privilege Register Failed,${SVCMGT_SVC} is not ready!"
            exit 30
        fi
        if [[ ${return_code} =~ ^[2-3] ]];then
            echo "${SVCMGT_SVC} is ready!"
            break
        else
            echo "${SVCMGT_SVC} is not ready,wait for 5s......"
            sleep 5
            continue
        fi
    done
}


## 2. 获取注册请求体
function obtain_register_body() {
    if [ ! -f "/opt/shellScripts/register_privilege_body.json" ];then
        echo "There is no register body file,register privilege failed!"
        exit 10
    fi
    jq -r . /opt/shellScripts/register_privilege_body.json
    if [ $? -ne "0" ];then
        echo "There is something wrong with the file of register_privilege_body.json,please check the format of the file!"
        exit 15
    fi
    #register_body=`cat /opt/shellScripts/register_privilege_body.json`
}

## 3. 注册服务
function registered_privilege() {

    virtAdminToken=$(curl -s -i -X POST \
        http://${KONG_SVC}:${KONG_PORT}/sys/auth/v1/tokens \
        -H 'Content-Type: application/json' \
        -d '{"auth":{"identity":{"methods":["password"],"password":{"user":{"name":"virtAdmin"},"password":"Passw0rd@_"}}}}' \
        | grep X-Subject-Token | sed 's/:/\n/g' | sed '1d')

    echo "virtAdminToken: ${virtAdminToken}"

    return_code=$(curl -o /dev/null -s -w %{http_code} -X POST \
        http://${SVCMGT_SVC}:${SVCMGT_PORT}/v1/services/${SERVICE_NAME}/privileges/init \
        -H 'Content-Type: application/json' \
        -H 'X-Auth-Token: '"${virtAdminToken}"'' \
        -d @/opt/shellScripts/register_privilege_body.json)
    echo "return_code: ${return_code}"
    if [[ ${return_code} =~ ^[2-3] ]];then
        echo "privilege register successfully!"
        register_statement
    else
        echo "Register privilege failed!"
        exit 20
    fi
}

function register_statement() {
    RESOURCE_KEY=`cat /opt/shellScripts/register_privilege_body.json | jq .cloudService.product_name | sed 's/"//g'`
    for (( times=0; times<= 5; times++ ))
    do
        virtAdminToken=$(curl -s -i -X POST \
            http://${KONG_SVC}:${KONG_PORT}/sys/auth/v1/tokens \
            -H 'Content-Type: application/json' \
            -d '{"auth":{"identity":{"methods":["password"],"password":{"user":{"name":"virtAdmin"},"password":"Passw0rd@_"}}}}' \
            | grep X-Subject-Token | sed 's/:/\n/g' | sed '1d')

        echo "virtAdminToken: ${virtAdminToken}"

        return_code=$(curl -o /dev/null -s -w %{http_code} -X GET \
           -H 'X-Auth-Token: '"${virtAdminToken}"'' \
           http://${SVCMGT_SVC}:${SVCMGT_PORT}/v1/service/statement?namespace=${NAMESPACE}\&servicename=${SERVICE_NAME}\&key=${RESOURCE_KEY}\&type=privilege)

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

function main() {
    check_service
    obtain_register_body
    registered_privilege
}

main
