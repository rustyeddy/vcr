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
echo "Turning video on"
mosquitto_pub -t camera/${baron} -m on
sleep 1
echo "Turning video off"
mosquitto_pub -t camera/${baron} -m on
echo "All done."
