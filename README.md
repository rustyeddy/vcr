# Redeye Smart Network Camera

RedEye is *smart* IP camera software built to run on inexpensive
computers with connected cameras like the Raspberry Pi with a CSI
camera. 

The idea is to be able to *control* a *network of cameras* providing
various video stream(s) to _Computer Vision_ algorithms.

## Working Features

RedEye was built to _stream video over IP_ as well take control
commands from a network client. With that in mind, RedEye was built
with the following features:

1. Device support is provided by OpenCV, releiving me of programming. 
2. Stream M-JPEG files over HTTP
7. Plugin system for configurable *Computer Vision Pipelines CVP*
3. REST API for the Config, Control and Storage Interfaces
4. MQTT messaging support for the above APIs
5. Websocket support for APIs
6. Embedded Webapp to control camera

## Near Term Roadmap

0. Off device video pipeline (stream to nano for CV Pipeline)
4. Improve the CV pipeline
1. Configurable cloud storage options
2. OpenCV to become a plugin
4. Stream only (M-JPEG) only support for cameras

## Supported Platforms

+ Raspberry Pi 3/4 + CSI Camera
+ Jetson Nano + CSI Camera
+ Ubuntu 19 Desktop + USB Cam (V4L)
+ Macbook Pro and Ait + Built in Camera
+ TODO Windows

### TODO OpenCV Plugin and Stream Only

+ Raspberry Pi Zero (stream only)
+ esp32 cam (st)

## OpenCV Plugin and Performance

RedEye is built with _OpenCV_ and hence takes advantage of the
powerful and flexible _device support_ provided by _OpenCV_. With
that, we get an amazing amount of power and flexibility right out of
the box, and do not have to do too much hard work to get there.

However, it does come at quite a footprint regarding memory, and the
build time on smaller devices is _ridiculous_ by todays standards (I
feel like a spoiled brat).

The idea then is to simply have the camera stream video to the A/I
module on another system. That requires the following to that going, 
Computer Vision module to read streaming video from network. 

That way, the smart module, can just suck the video down from a player
that only knows how to stream the video.


## OpenCV and Pipeline Plugins

+ Built with OpenCV
+ Video Pipeline plugins
  + Face detection

## APIs

### Camera Control

- Play
- Pause 
- Snap

### Camera Config

- Resolution
- Frames Per Second
- Format

### Storage

- Location
- GetClips
- GetSnapshots
- SaveClip
- SaveSnapshot

## Otto Discovery

+ Otto Discovery with MQTT
  + requires an MQTT broker
  + optional if broker is MQTT broker is NOT present


See Todo.org
