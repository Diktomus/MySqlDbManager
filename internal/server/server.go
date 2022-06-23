package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router  *mux.Router
	address string
}

func NewServer(router *mux.Router, address string) *Server {
	return &Server{
		router:  router,
		address: address,
	}
}

func (server *Server) Run() error {
	return http.ListenAndServe(server.address, server.router)
}
