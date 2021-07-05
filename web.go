package redeye

import (
	"encoding/json"
	"net/http"
	"sync"
)

type WebServer struct {
	Addr     string
	Basepath string // /redeye
	Handlers []string
}

var (
	web    WebServer
	stream Stream

	successResponse = map[string]string{"success": "true"}
	errorResponse   = map[string]string{"success": "false"}
)

func GetWebServer(Addr, Path string) (web *WebServer) {
	web = &WebServer{
		Addr:     Addr,
		Basepath: Path,
		Handlers: nil,
	}

	web.RegisterHandlerFunc(Path+"/health", health)
	web.RegisterHandlerFunc(Path+"/cameras", GetCameras)
	web.RegisterHandler("/mjpeg", stream)
	web.RegisterHandler("/ws", WSServer{})
	return web
}

func (web *WebServer) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	http.ListenAndServe(web.Addr, nil)
}

// Register a REST handler to match the given base URL.
func (web *WebServer) RegisterHandler(path string, handler http.Handler) {
	web.Handlers = append(web.Handlers, path)
	http.Handle(path, handler)
}

// Register a REST handler to match the given base URL.
func (web *WebServer) RegisterHandlerFunc(path string, handler func(http.ResponseWriter, *http.Request)) {
	web.Handlers = append(web.Handlers, path)
	http.HandleFunc(path, handler)
}

func health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(successResponse)
}
