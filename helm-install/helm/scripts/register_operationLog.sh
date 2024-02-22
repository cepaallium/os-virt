#!/bin/bash

## 1. 检查服务是否就绪
function check_service() {
    for (( times=0; times<= 30; times++ ))
    do
        return_code=$(curl -o /dev/null -s -w %{http_code} http://${CILL_SVC}:${CILL_PORT}/v1/rest_health)
        if [ ${times} -eq 30 ];then
            echo "Menu Register Failed,${CILL_SVC} is not ready!"
            exit 30
        fi
        if [[ ${return_code} =~ ^[2-3] ]];then
            echo "${CILL_SVC} is ready"
            break
        else
            echo "${CILL_SVC} is not ready,wait for 5s......"
            sleep 5
            continue
        fi
    done
}


## 2. 获取注册请求体
function obtain_register_body() {
    if [ ! -f "/opt/shellScripts/register_operationLog_body.json" ];then
        echo "There is no register body file,register operationLog failed!"
        exit 10
    fi
    jq -r . /opt/shellScripts/register_operationLog_body.json
    if [ $? -ne "0" ];then
        echo "There is something wrong with the file of register_operationLog_body.json,please check the format of the file!"
        exit 15
    fi
    #register_body=`cat /opt/shellScripts/register_operationLog_body.json`
}

## 3. 注册服务
function registered_operationLog() {
    cillAdminToken=$(curl -s -i -X POST \
        http://${KONG_SVC}:${KONG_PORT}/sys/auth/v1/tokens \
        -H 'Content-Type: application/json' \
        -d '{"auth":{"identity":{"methods":["password"],"password":{"user":{"name":"cillAdmin"},"password":"Passw0rd@_"}}}}' \
        | grep X-Subject-Token | sed 's/:/\n/g' | sed '1d')
    echo "cillAdminToken: ${cillAdminToken}"
    return_code=$(curl -o /dev/null -s -w %{http_code} -X POST \
        http://${KONG_SVC}:${KONG_PORT}/sys/oapi/v1/i18n/data \
        -H 'Content-Type: application/json' \
        -H 'X-Auth-Token: '"${cillAdminToken}"'' \
        -d @/opt/shellScripts/register_operationLog_body.json)
        echo "return_code: ${return_code}"
    if [[ ${return_code} =~ ^[2-3] ]];then
        echo "OperationLog register successfully!"
    else
        echo "Register OperationLog failed!"
        exit 20
    fi
}

function main() {
    check_service
    obtain_register_body
    registered_operationLog
}

main
