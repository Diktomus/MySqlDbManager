package main

import (
	"fmt"
	"github/mysql-dbmanager/internal/config"
	DbController "github/mysql-dbmanager/internal/controller"
	"github/mysql-dbmanager/internal/delivery"
	"net/http"

	"database/sql"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbConfig, err := config.Init()
	if err != nil {
		return
	}

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?&charset=utf8&interpolateParams=true", dbConfig.Login, dbConfig.Password, dbConfig.Ip, dbConfig.Port, dbConfig.DbName)

	database, _ := sql.Open("mysql", dataSourceName)

	database.SetMaxOpenConns(dbConfig.MaxConns)

	err = database.Ping()
	if err != nil {
		return
	}

	controller := DbController.NewController(database)
	controller.Init()

	handlers := &delivery.Handlers{
		Controller: controller,
	}

	gorillaMux := mux.NewRouter()

	gorillaMux.HandleFunc("/", handlers.GetTablesHandler).Methods("GET")
	gorillaMux.HandleFunc("/{table}", handlers.GetEntriesHandler).Methods("GET")
	gorillaMux.HandleFunc("/{table}/{id}", handlers.GetEntryHandler).Methods("GET")
	gorillaMux.HandleFunc("/{table}", handlers.CreateEntryHandler).Methods("PUT")
	gorillaMux.HandleFunc("/{table}/{id}", handlers.UpdateEntryHandler).Methods("POST")
	gorillaMux.HandleFunc("/{table}/{id}", handlers.DeleteEntryHandler).Methods("DELETE")

	fmt.Printf("Start http server on %s:%d\n", dbConfig.Ip, dbConfig.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", dbConfig.Ip, dbConfig.Port), gorillaMux)
}
