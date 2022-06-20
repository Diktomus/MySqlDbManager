package server

import (
	"fmt"
	"github/mysql-dbmanager/internal/config"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	config *config.MySqlDbConfig
}

func NewServer(router *mux.Router, config *config.MySqlDbConfig) *Server {
	return &Server{
		router: router,
		config: config,
	}
}

func (server *Server) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", 8080), server.router)
}
