package redeye

import (
	"encoding/json"
	"net/http"
)

var (
	cameras map[string]*Camera
)

func init() {
	cameras = make(map[string]*Camera)
}

type Camera struct {
	Name string
	Addrport string
}

func NewCamera(name string) *Camera {
	cam := &Camera{Name: name, Addrport: name}
	cameras[name]  = cam
	return cam;
}

func GetCameras(w http.ResponseWriter, r *http.Request) {
	clist := GetCameraList()
	json.NewEncoder(w).Encode(clist)
}

func GetCameraList() (clist []*Camera) {
	for _, cam := range cameras {
		clist = append(clist, cam)
	}
	return clist
}


func (cam *Camera) Handler(w http.ResponseWriter, req *http.Request) {
	GetCameras(w, req)
}

