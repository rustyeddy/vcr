package redeye

import (
	"log"
	"sync"
	"encoding/json"
	"net/http"
)

type WebServer struct {
	Addr string
	Basepath string				// /redeye
	Handlers []string
}

var (
	successResponse = map[string]string{"success": "true"}
	errorResponse   = map[string]string{"success": "false"}
)

func NewWebServer(Addr, Path string) (web *WebServer) {
	web = &WebServer{
		Addr: Addr,
		Basepath: Path,
		Handlers: nil,
	}
	http.HandleFunc(Path + "/health", health)
	return web
}

func (web *WebServer) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	
	log.Printf("New HTTP Server %s created", web.Addr)
	http.ListenAndServe(web.Addr, nil)
}

// Register a REST handler to match the given base URL.
func (web *WebServer) RegisterHandler(path string, handler http.Handler) {
	log.Println("WebServer ~ Adding handler: ", path)
	web.Handlers = append(web.Handlers, path)
	http.Handle(path, handler)
}

// Register a REST handler to match the given base URL.
func (web *WebServer) RegisterHandlerFunc(path string, handler func(http.ResponseWriter, *http.Request)) {
	log.Println("WebServer ~ Adding handler: ", path)
	web.Handlers = append(web.Handlers, path)
	http.HandleFunc(path, handler)
}

func health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(successResponse)
}
