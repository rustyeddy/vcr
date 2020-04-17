#!/bin/bash
curl http://localhost:8888/health
sleep 2
curl http://localhost:8888/config
sleep 2
curl http://localhost:8888/messanger
sleep 2
