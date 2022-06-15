package delivery

import (
	"fmt"
	"github/mysql-dbmanager/internal/adapter"
	"github/mysql-dbmanager/internal/controller"
	"github/mysql-dbmanager/internal/utils"
	"net/http"
	"strconv"
)

type Handlers struct {
	Controller *controller.Controller
}

func (h *Handlers) GetTablesHandler(resp http.ResponseWriter, req *http.Request) {

	fmt.Fprintf(resp, "Show tables names:\n")

	for _, table := range h.Controller.Tables {
		fmt.Fprintf(resp, "%s\n", table)
	}

	fmt.Println("called GetTablesHandler")
}

func (h *Handlers) GetEntriesHandler(resp http.ResponseWriter, req *http.Request) {
	tableName := utils.GetVariable("table", req)

	limit, err := strconv.Atoi(req.FormValue("limit"))
	if err != nil {
		limit = 10
	}
	offset, _ := strconv.Atoi(req.FormValue("offset"))

	rows := h.Controller.GetRows(tableName, limit, offset)
	defer rows.Close()

	adaptedRows, err := adapter.AdaptSqlRows(rows)
	if err != nil {
		return
	}
	adapter.WriteRowsToResp(adaptedRows, resp)

	fmt.Println("called GetEntriesHandler")
}

func (h *Handlers) GetEntryHandler(resp http.ResponseWriter, req *http.Request) {
	tableName := utils.GetVariable("table", req)
	id, _ := strconv.Atoi(utils.GetVariable("id", req))

	row := h.Controller.GetRow(tableName, id)

	columns := h.Controller.GetColumns(tableName)

	adaptedRow, err := adapter.AdaptSqlRow(row, columns)
	if err != nil {
		return
	}
	fmt.Fprintf(resp, "Row with id = %d\n", id)
	adapter.WriteRowsToResp([]adapter.AdaptedRow{adaptedRow}, resp)

	fmt.Println("called GetEntryHandler")
}

func (h *Handlers) CreateEntryHandler(resp http.ResponseWriter, req *http.Request) {
	tableName := utils.GetVariable("table", req)

	req.ParseMultipartForm(32 << 20)

	result := h.Controller.CreateRow(tableName, req.Form)

	utils.WriteResultToResp(result, resp)

	fmt.Println("called CreateEntryHandler")
}

func (h *Handlers) UpdateEntryHandler(resp http.ResponseWriter, req *http.Request) {
	tableName := utils.GetVariable("table", req)
	id, _ := strconv.Atoi(utils.GetVariable("id", req))

	req.ParseMultipartForm(32 << 20)

	result := h.Controller.UpdateRow(tableName, req.Form, id)

	utils.WriteResultToResp(result, resp)

	fmt.Println("called UpdateEntryHandler")
}

func (h *Handlers) DeleteEntryHandler(resp http.ResponseWriter, req *http.Request) {
	tableName := utils.GetVariable("table", req)
	id, _ := strconv.Atoi(utils.GetVariable("id", req))

	h.Controller.DeleteRow(tableName, id)

	fmt.Fprintf(resp, "Row with id = %d was deleted\n", id)

	fmt.Println("called DeleteEntryHandler")
}
