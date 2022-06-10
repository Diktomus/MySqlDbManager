package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type AdaptedRow map[string]interface{}

type ScanBuffer struct {
	Pointers []interface{}
	Values   []interface{}
}

func GetVariable(variableName string, req *http.Request) string {
	vars := mux.Vars(req)
	variable := vars[variableName]

	return variable
}

func PrepareScanBuffer(bufferSize int) ScanBuffer {
	scanBuffer := ScanBuffer{
		Pointers: make([]interface{}, bufferSize),
		Values:   make([]interface{}, bufferSize),
	}

	for i := 0; i < len(scanBuffer.Values); i++ {
		scanBuffer.Pointers[i] = &scanBuffer.Values[i]
	}

	return scanBuffer
}

func ResolveAdaptedRow(columns []string, values []interface{}) AdaptedRow {
	adaptedRow := make(AdaptedRow, 0)
	for i, column := range columns {
		if value, ok := (values[i]).([]byte); ok {
			adaptedRow[column] = string(value)
		}
	}
	return adaptedRow
}

func ParseStr(s string) interface{} {
	result, err := strconv.ParseFloat(s, 2)
	if err != nil {
		return s
	}
	return result
}

func AdaptDbRowsData(rows *sql.Rows) []AdaptedRow {
	adaptedRows := make([]AdaptedRow, 0)
	columns, _ := rows.Columns()
	scanBuffer := PrepareScanBuffer(len(columns))

	for rows.Next() {
		err := rows.Scan(scanBuffer.Pointers...)
		if err != nil {
			panic(err)
		}
		adaptedRow := ResolveAdaptedRow(columns, scanBuffer.Values)
		adaptedRows = append(adaptedRows, adaptedRow)
	}
	return adaptedRows
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
