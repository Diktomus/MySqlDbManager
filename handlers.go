package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
)

type Handlers struct {
	DB         *sql.DB
	controller *Controller
}

func (h *Handlers) GetTablesHandler(resp http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(resp, "Show tables names:\n")

	for _, table := range h.controller.Tables {
		fmt.Fprintf(resp, "%s\n", table)
	}

	fmt.Println("called GetTablesHandler")
}

func (h *Handlers) GetEntriesHandler(resp http.ResponseWriter, req *http.Request) {
	tableName := GetVariable("table", req)

	limit, err := strconv.Atoi(req.FormValue("limit"))
	if err != nil {
		limit = 10
	}
	offset, _ := strconv.Atoi(req.FormValue("offset"))

	rows := h.controller.GetRows(tableName, limit, offset)

	adaptedRows := AdaptDbRowsData(rows)
	rows.Close()

	WriteRowsToResp(adaptedRows, resp)

	fmt.Println("called GetEntriesHandler")
}

func (h *Handlers) GetEntryHandler(resp http.ResponseWriter, req *http.Request) {
	tableName := GetVariable("table", req)
	id, _ := strconv.Atoi(GetVariable("id", req))

	row := h.controller.GetRow(tableName, id)

	columns := h.controller.GetColumns(tableName)

	scanBuffer := PrepareScanBuffer(len(columns))

	err := row.Scan(scanBuffer.Pointers...)
	if err == sql.ErrNoRows {
		http.NotFound(resp, req)
		return
	} else if err != nil {
		http.Error(resp, http.StatusText(500), 500)
	}

	adaptedRow := ResolveAdaptedRow(columns, scanBuffer.Values)

	fmt.Fprintf(resp, "Row with id = %d\n", id)
	WriteRowsToResp([]AdaptedRow{adaptedRow}, resp)

	fmt.Println("called GetEntryHandler")
}

func (h *Handlers) CreateEntryHandler(resp http.ResponseWriter, req *http.Request) {
	tableName := GetVariable("table", req)

	req.ParseMultipartForm(32 << 20)

	result := h.controller.CreateRow(tableName, req.Form)

	WriteResultToResp(result, resp)

	fmt.Println("called CreateEntryHandler")
}

func (h *Handlers) UpdateEntryHandler(resp http.ResponseWriter, req *http.Request) {
	tableName := GetVariable("table", req)
	id, _ := strconv.Atoi(GetVariable("id", req))

	req.ParseMultipartForm(32 << 20)

	result := h.controller.UpdateRow(tableName, req.Form, id)

	WriteResultToResp(result, resp)

	fmt.Println("called UpdateEntryHandler")
}

func (h *Handlers) DeleteEntryHandler(resp http.ResponseWriter, req *http.Request) {
	tableName := GetVariable("table", req)
	id, _ := strconv.Atoi(GetVariable("id", req))

	h.controller.DeleteRow(tableName, id)

	fmt.Fprintf(resp, "Row with id = %d was deleted\n", id)

	fmt.Println("called DeleteEntryHandler")
}
