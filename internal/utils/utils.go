package utils

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func GetVariable(variableName string, req *http.Request) string {
	vars := mux.Vars(req)
	variable := vars[variableName]

	return variable
}

func GetDBName(dbSource string) string {
	slashPos := strings.Index(dbSource, "/")
	questionPos := strings.Index(dbSource, "?")
	if questionPos == -1 {
		return dbSource[slashPos+1:]
	}
	return dbSource[slashPos:questionPos]
}

func ParseStr(s string) interface{} {
	result, err := strconv.ParseFloat(s, 2)
	if err != nil {
		return s
	}
	return result
}

func WriteResultToResp(result sql.Result, resp http.ResponseWriter) {
	lastInsertedId, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()

	fmt.Fprintf(resp, "Last inserted id: %d\n", lastInsertedId)
	fmt.Fprintf(resp, "Rows Affected: %d\n", rowsAffected)
}
