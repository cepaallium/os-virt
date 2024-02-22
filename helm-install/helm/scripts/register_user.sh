#!/bin/sh

## 1. 检查服务是否就绪
function check_service() {
    for (( times=0; times<= 30; times++ ))
    do
        return_code=$(curl -o /dev/null -s -w %{http_code} http://${SVCMGT_SVC}:${SVCMGT_PORT}/v1/rest_health)
        if [ ${times} -eq 30 ];then
            echo "Billing Register Failed,${SVCMGT_SVC} is not ready!"
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
    if [ ! -f "/opt/shellScripts/register_user_body.json" ];then
        echo "There is no register body file,register user failed!"
        exit 10
    fi
    jq -r . /opt/shellScripts/register_user_body.json
    if [ $? -ne "0" ];then
        echo "There is something wrong with the file of register_user_body.json,please check the format of the file!"
        exit 15
    fi
}


## 3. 注册服务
function registered_user() {
  X_AUTH_TOKEN=$(curl -X POST  http://${BOLT_SVC}:${BOLT_PORT}/v1/tokens  -d "{\"auth\":{\"identity\":{\"methods\":[\"password\"],\"password\":{\"user\":{\"name\":\"cillAdmin\"},\"password\":\"Passw0rd@_\"}}}}" -H "Content-Type: application/json" -i | grep "X-Subject-Token" | cut -d " "  -f 2)
    return_code=$(curl -o /dev/null -s -w %{http_code} -X POST \
        http://${SVCMGT_SVC}:${SVCMGT_PORT}/v1/users \
        -H "Content-Type: application/json" \
        -H "X-Auth-Token: ${X_AUTH_TOKEN}" \
        -d @/opt/shellScripts/register_user_body.json)
    echo "return_code: ${return_code}"
    if [[ ${return_code} =~ ^[2-3] ]];then
        echo "user register successfully!"
    else
        echo "Register user failed!"
        exit 20
    fi
}

function main() {

    check_service

    obtain_register_body
    registered_user

}

main