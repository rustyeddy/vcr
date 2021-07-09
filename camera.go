package redeye

import (
	"fmt"
	"encoding/json"
)

type Camera struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
	Port int	`json:"port"`
	URI	 string `json:"uri"`
}

func NewCamera(camstr string) *Camera {

	fmt.Println("Camstr: ", camstr)

	var cam Camera
	err := json.Unmarshal([]byte(camstr), &cam)
	if err != nil {
		fmt.Println("ERROR - unmarshal camera json", err)
		return nil
	}

	//cam := &Camera{Name: name, Addr: name, Port: 8080}
	cameras[cam.Name] = &cam
	return &cam;
}

