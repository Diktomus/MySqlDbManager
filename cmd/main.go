package main

import (
	"fmt"
	"github/mysql-dbmanager/internal/config"
	DbController "github/mysql-dbmanager/internal/controller"
	"github/mysql-dbmanager/internal/delivery"
	"github/mysql-dbmanager/internal/server"
	"github/mysql-dbmanager/internal/utils"
	"os"

	"database/sql"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	dbConfig, err := config.Load("../config/app")
	if err != nil {
		log.Error().Err(err).Msg("config.Load")
		return
	}
	database, err := sql.Open(dbConfig.DBDriver, dbConfig.DBSource)
	defer database.Close()

	database.SetMaxOpenConns(dbConfig.MaxConns)
	err = database.Ping()
	if err != nil {
		log.Error().Err(err).Msg("Ping database")
		return
	}
	dbName := utils.GetDBName(dbConfig.DBSource)
	controller := DbController.NewController(database, dbName)
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

	message := fmt.Sprintf("Starting http server on %s\n", dbConfig.ServerAddress)
	log.Info().Msg(message)
	server := server.NewServer(router, dbConfig.ServerAddress)
	err = server.Run()
	if err != nil {
		log.Error().Err(err).Msg("server.Run")
	}
}
