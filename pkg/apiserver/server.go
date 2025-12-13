package apiserver

import (
	"log"
	"net/http"
)

type Server struct {
	Addr    string
	Handler http.Handler
}

func (s *Server) Run() {
	log.Println("API server listening on", s.Addr)
	log.Fatal(http.ListenAndServe(s.Addr, s.Handler))
}
