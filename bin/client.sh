#!/bin/bash
camera=baron

echo -n "Getting health .. "
curl http://localhost:8000/health
sleep 1
echo -n "Getting config .. "
curl http://localhost:8000/config
sleep 1
echo -n "Getting messanger .. "
curl http://localhost:8000/messanger
sleep 1
echo "Turning video on for 10 seconds .. "
mosquitto_pub -t camera/${camera} -m on
sleep 10
echo "Turning video off"
mosquitto_pub -t camera/${camera} -m off
echo "All done."
