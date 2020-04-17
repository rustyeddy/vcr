#!/bin/bash
curl http://localhost:8888/health
sleep 1
curl http://localhost:8888/config
sleep 1
curl http://localhost:8888/messanger
sleep 1
mosquitto_pub -t camera/control -m on
