package main

import (
	"flag"
	"fmt"
	"net/http"

	"database/sql"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	login := flag.String("login", "", "database login")
	passwd := flag.String("passwd", "", "database password")
	ip := flag.String("ip", "localhost", "database ip")
	port := flag.String("port", "", "database port")
	dbName := flag.String("db_name", "", "database name")
	maxConns := flag.Int("max_conns", 5, "max database connections")

	flag.Parse()

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?&charset=utf8&interpolateParams=true", *login, *passwd, *ip, *port, *dbName)

	database, _ := sql.Open("mysql", dataSourceName)

	database.SetMaxOpenConns(*maxConns)

	err := database.Ping()
	if err != nil {
		panic(err)
	}

	controller := NewController(database)

	handlers := &Handlers{
		DB:         database,
		controller: controller,
	}

	gorillaMux := mux.NewRouter()

	gorillaMux.HandleFunc("/", handlers.GetTablesHandler).Methods("GET")
	gorillaMux.HandleFunc("/{table}", handlers.GetEntriesHandler).Methods("GET")
	gorillaMux.HandleFunc("/{table}/{id}", handlers.GetEntryHandler).Methods("GET")
	gorillaMux.HandleFunc("/{table}", handlers.CreateEntryHandler).Methods("PUT")
	gorillaMux.HandleFunc("/{table}/{id}", handlers.UpdateEntryHandler).Methods("POST")
	gorillaMux.HandleFunc("/{table}/{id}", handlers.DeleteEntryHandler).Methods("DELETE")

	fmt.Printf("Start http server on %s:%s\n", *ip, *port)
	http.ListenAndServe(fmt.Sprintf("%s:%s", *ip, *port), gorillaMux)
}
