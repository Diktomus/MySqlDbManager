package utils

import (
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
