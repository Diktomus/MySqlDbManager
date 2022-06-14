package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetVariable(variableName string, req *http.Request) string {
	vars := mux.Vars(req)
	variable := vars[variableName]

	return variable
}

func ParseStr(s string) interface{} {
	result, err := strconv.ParseFloat(s, 2)
	if err != nil {
		return s
	}
	return result
}

func WriteRowsToResp(adaptedRows []AdaptedRow, resp http.ResponseWriter) {
	for _, adaptedRow := range adaptedRows {
		fmt.Fprintf(resp, "%+v\n", adaptedRow)
	}
}

func WriteResultToResp(result sql.Result, resp http.ResponseWriter) {
	lastInsertedId, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()

	fmt.Fprintf(resp, "Last inserted id: %d\n", lastInsertedId)
	fmt.Fprintf(resp, "Rows Affected: %d\n", rowsAffected)
}
