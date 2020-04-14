# Redeye Smart Network Camera

RedEye is *smart* IP camera software that runs on lots of inexpensive
computers and connected cameras, pretty much anythin OpenCV runs
on. RedEye uses OpenCV _VideoCapture_ device to read anything from a
_network IP camera_ to an _mp4 video file_.

## Overview 

I built _redeye_ to control cheap cameras attached to cheap
computers, specifically the _Raspberry Pi and CSI camera_.


## Open Source and Inexpensive Hardware

RedEye was written, originally to control a Raspberry Pi with a
connected CSI camera, over a network. I used GOCV to build the camera
software, by virtue of _OpenCV_ VideoCapture() Device it was trivial
to support other platforms, such as follows:

### Camera with OpenCV Support

All RedEye software is built with OpenCV, which is a very large and
powerful collection of software, quite a lot of resources go into 
bit goes into building it and supportnig 

RedEye was designed to use _plugins_ allowing a specific camera to
only load the additional functionality it requires for a particular
application. 

This is important 


+ Raspberry Pi 3/4 + CSI Camera
+ Jetson Nano + CSI Camera
+ Ubuntu 19 Desktop + USB Cam (V4L)
+ Macbook Pro and Ait + Built in Camera
+ TODO Raspberry Pi Zero (stream only)



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
