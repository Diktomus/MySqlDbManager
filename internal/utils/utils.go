package utils

import (
	"net/http"

	"github.com/gorilla/mux"
)

func GetVariable(variableName string, req *http.Request) string {
	vars := mux.Vars(req)
	variable := vars[variableName]

	return variable
}
