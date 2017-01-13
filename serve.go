package main

import (
	"net/http"
)

var httpListenAndServe = http.ListenAndServe
var httpWriterSetContentType = func(w http.ResponseWriter, value string) {
	w.Header().Set("Content-Type", value)
}

type Server interface {
	Run()
}

type Serve struct {
	Service Servicer
}

func (m *Serve) Run() error {
	if err := httpListenAndServe(":8080", m); err != nil {
		return err
	}
	return nil
}

func (m *Serve) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	httpWriterSetContentType(w, "application/json")
	switch req.URL.Path {
	case "/v1/docker-flow-swarm-listener/notify-services":
		services, _ := m.Service.GetServices()
		go m.Service.NotifyServicesCreate(services, 10, 5)
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func NewServe(service Servicer) *Serve {
	return &Serve{
		Service: service,
	}
}
