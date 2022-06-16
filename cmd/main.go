package main

import (
	"fmt"
	"github/mysql-dbmanager/internal/config"
	DbController "github/mysql-dbmanager/internal/controller"
	"github/mysql-dbmanager/internal/delivery"
	"github/mysql-dbmanager/internal/server"
	"os"

	"database/sql"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	dbConfig, err := config.Init()
	if err != nil {
		log.Error().Err(err).Msg("config.Init")
		return
	}

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?&charset=utf8&interpolateParams=true", dbConfig.Login, dbConfig.Password, dbConfig.Ip, dbConfig.Port, dbConfig.DbName)

	database, err := sql.Open("mysql", dataSourceName)
	database.SetMaxOpenConns(dbConfig.MaxConns)
	err = database.Ping()
	if err != nil {
		log.Error().Err(err).Msg("Ping database")
		return
	}

	controller := DbController.NewController(database, dbConfig.DbName)
	err = controller.Init()
	if err != nil {
		return
	}

	handlers := &delivery.Handlers{
		Controller: controller,
	}

	router := mux.NewRouter()
	router.HandleFunc("/", handlers.GetTablesHandler).Methods("GET")
	router.HandleFunc("/{table}", handlers.GetEntriesHandler).Methods("GET")
	router.HandleFunc("/{table}/{id}", handlers.GetEntryHandler).Methods("GET")
	router.HandleFunc("/{table}", handlers.CreateEntryHandler).Methods("PUT")
	router.HandleFunc("/{table}/{id}", handlers.UpdateEntryHandler).Methods("POST")
	router.HandleFunc("/{table}/{id}", handlers.DeleteEntryHandler).Methods("DELETE")

	server := server.NewServer(router, dbConfig)
	server.Run()
}
