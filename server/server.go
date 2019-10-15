package server

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server interface {
	MapRoutes([]*Route)
	Start()
}

type MuxServer struct {
	Port      string
	Listening bool
	router    *mux.Router
}

type Route struct {
	Endpoint string
	Handler  func(http.ResponseWriter, *http.Request)
	Methods  []string
}

func NewMuxServer(port string) *MuxServer {
	s := new(MuxServer)
	s.Port = port
	return s
}

func (s *MuxServer) MapRoutes(routeMap []*Route) {
	router := mux.NewRouter()
	for _, route := range routeMap {
		router.HandleFunc(route.Endpoint, route.Handler).Methods(route.Methods...)
	}
	s.router = router
}

func (s *MuxServer) Start() {
	// Start listening using the mux router
	log.Println("Listening on port " + s.Port)
	http.ListenAndServe(":"+s.Port, s.router)
}
