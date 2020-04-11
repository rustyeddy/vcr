/*

RedEye the smart camera software.

- MQTT Play and Pause Video

- GET		/pi/health
- GET		/api/config
- POST|PUT	/api/config/?key=val&key=val

*/
package main

// The Camera Status
type CameraStatus struct {
	Name         string
	Addr         string
	Status       string
	PipelineName string
}
