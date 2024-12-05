#!/bin/sh
TOKEN=`curl -X POST "http://wazigate-dashboard.staging.waziup.io/auth/token" -H  "accept: text/plain" -H  "Content-Type: application/json" -d "{\"username\":\"admin\",\"password\":\"loragateway\"}"`

curl -X POST "http://wazigate-dashboard.staging.waziup.io/apps" -H "accept: */*" -H "Content-Type: application/json;charset=utf-8" -H "Authorization:Bearer $TOKEN" -d '"waziupiot/majiup:latest"'

curl -X POST "http://wazigate-dashboard.staging.waziup.io/apps/waziupiot.majiup" -H "accept: */*" -H "Content-Type: application/json;charset=utf-8" -H "Authorization:Bearer $TOKEN" -d '{"action":"start"}'
