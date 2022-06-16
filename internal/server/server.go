package server

import (
	"fmt"
	"github/mysql-dbmanager/internal/config"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
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
	message := fmt.Sprintf("Start http server on %s:%d\n", server.config.Ip, server.config.Port)
	log.Info().Msg(message)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", server.config.Ip, server.config.Port), server.router)
}
