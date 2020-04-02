/*

RedEye the smart camera software.

- MQTT Play and Pause Video

- GET		/pi/health
- GET		/api/config
- POST|PUT	/api/config/?key=val&key=val

*/
package main

// Camera representation to outside world
type Camera struct {
	Name      string
	Location  string
	Pipeline  string
	Recording bool

	Clips []string
}
