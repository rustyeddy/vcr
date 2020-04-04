# redeye

RedEye is *smart* camera software built with Open Computer Vision
(OpenCV) and with _simple control and configuration_ APIs.

## What Works

For list of what works and what is being worked on is in the TODO
lsit. A brief of what currently works:

### M-JPEG Stream Server

Stream MJpeg from most OSs with video cam, the following have been
tested at one point in time:

+ Raspberry Pi w/CSI 
+ Jetson Nano w/CSI
+ Desktop Ubuntu w/USB
+ MacOS laptop (air and pro) with builtin Cam
- Windows (TBD)

### OpenCV and Pipeline Plugins

+ Built with OpenCV
+ Video Pipeline plugins
  + Face detection

### Remote Controls

+ Play/Pause over MQTT and Websockets
+ Configure Get configuration with REST

## Otto Discovery

+ Otto Discovery with MQTT
  + requires an MQTT broker
  + optional if broker is MQTT broker is NOT present


See Todo.org
