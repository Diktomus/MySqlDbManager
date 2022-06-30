package utils

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

func GetVariable(variableName string, req *http.Request) string {
	vars := mux.Vars(req)
	variable := vars[variableName]

	return variable
}

func ParseUrlValues(urlValues url.Values) (columns []string, values []interface{}, err error) {
	columns = make([]string, 0)
	values = make([]interface{}, 0)
	for column, columnValues := range urlValues {
		columns = append(columns, column)
		if len(columnValues) > 1 || len(columnValues) == 0 {
			return nil, nil, &ErrWrongUrlValues{valueName: column}
		} else {
			value, err := strconv.Atoi(columnValues[0])
			if err != nil {
				values = append(values, columnValues[0])
			} else {
				values = append(values, value)
			}
		}
	}
	return columns, values, nil
}
