package redeye

import (
	"encoding/json"
	"net/http"
)

// Service defines the interface the redeye services. Services
// is meant to help compose the concrete service.
type Service interface {
	Name() string
	Addr() string

	Config(cfg map[string]string)
	Start(cmdQ chan TLV) chan TLV
	Stop()

	GetConfig(w *http.ResponseWriter)
}

// ServiceMeta defines the meta data required by every service
type ServiceMeta struct {
	name   string
	addr   string
	config *Settings
	status *Settings
}

// NewService creates a new generic service meant to help compose
// a concrete service, like the rest or mqtt services.
func NewService(name, addr string) *ServiceMeta {
	return &ServiceMeta{name, addr, nil, &Settings{}}
}

// Name returns the name of the service
func (s *ServiceMeta) Name() string {
	return s.name
}

// Name returns the name of the service
func (s *ServiceMeta) Addr() string {
	return s.addr
}

// Config allows the caller to provide this service with specific
// configuration info, like source address and port number.
func (s *ServiceMeta) Config(cfg map[string]string) {
	// TODO be smart about replaceing, well just replace for now
	s.config = NewSettings(cfg)
}

//
func (s *ServiceMeta) GetConfig(w http.ResponseWriter) {
	panic(s.status != nil)
	json.NewEncoder(w).Encode(s.status)

}

func (s *ServiceMeta) Start(cmdQ chan TLV) (q chan TLV) {
	panic("Start method must be defined in the concrete class")
}
func (s *ServiceMeta) Stop() {}
