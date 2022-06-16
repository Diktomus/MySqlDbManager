package delivery

import (
	"fmt"
	"github/mysql-dbmanager/internal/adapter"
	"github/mysql-dbmanager/internal/controller"
	"github/mysql-dbmanager/internal/utils"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

type Handlers struct {
	Controller *controller.Controller
}

func (h *Handlers) GetTablesHandler(resp http.ResponseWriter, req *http.Request) {
	log.Info().Msg("Called GetTablesHandler")
	fmt.Fprintf(resp, "Show tables names:\n")

	for _, table := range h.Controller.Tables {
		fmt.Fprintf(resp, "%s\n", table)
	}
}

func (h *Handlers) GetEntriesHandler(resp http.ResponseWriter, req *http.Request) {
	log.Info().Msg("Called GetEntriesHandler")
	tableName := utils.GetVariable("table", req)

	limit, err := strconv.Atoi(req.FormValue("limit"))
	if err != nil {
		limit = 10
		log.Info().Str("limit", strconv.Itoa(limit)).Msg("GetEntriesHandler: param limit isn't set")
	}
	offset, err := strconv.Atoi(req.FormValue("offset"))
	if err != nil {
		offset = 1
		log.Info().Str("offset", strconv.Itoa(offset)).Msg("GetEntriesHandler: param offset isn't set")
	}

	rows, err := h.Controller.GetRows(tableName, limit, offset)
	defer rows.Close()
	if err != nil {
		log.Error().Err(err).Msg("")
		http.NotFound(resp, req)
		return
	}

	adaptedRows, err := adapter.AdaptSqlRows(rows)
	if err != nil {
		log.Error().Err(err).Msg("")
		http.NotFound(resp, req)
		return
	}
	adapter.WriteRowsToResp(adaptedRows, resp)
}

func (h *Handlers) GetEntryHandler(resp http.ResponseWriter, req *http.Request) {
	log.Info().Msg("Called GetEntryHandler")
	tableName := utils.GetVariable("table", req)
	id, err := strconv.Atoi(utils.GetVariable("id", req))
	if err != nil {
		log.Error().Err(err).Msg("")
		http.NotFound(resp, req)
		return
	}

	row := h.Controller.GetRow(tableName, id)

	columns := h.Controller.GetColumns(tableName)

	adaptedRow, err := adapter.AdaptSqlRow(row, columns)
	if err != nil {
		log.Error().Err(err).Msg("")
		http.NotFound(resp, req)
		return
	}
	fmt.Fprintf(resp, "Row with id = %d\n", id)
	adapter.WriteRowsToResp([]adapter.AdaptedRow{adaptedRow}, resp)
}

func (h *Handlers) CreateEntryHandler(resp http.ResponseWriter, req *http.Request) {
	log.Info().Msg("Called CreateEntryHandler")
	tableName := utils.GetVariable("table", req)

	req.ParseMultipartForm(32 << 20)

	result, err := h.Controller.CreateRow(tableName, req.Form)
	if err != nil {
		log.Error().Err(err).Msg("CreateEntryHandler")
		http.Error(resp, http.StatusText(500), 500)
		return
	}
	utils.WriteResultToResp(result, resp)
}

func (h *Handlers) UpdateEntryHandler(resp http.ResponseWriter, req *http.Request) {
	log.Info().Msg("Called UpdateEntryHandler")
	tableName := utils.GetVariable("table", req)
	id, err := strconv.Atoi(utils.GetVariable("id", req))
	if err != nil {
		log.Error().Err(err).Msg("")
		http.NotFound(resp, req)
		return
	}

	req.ParseMultipartForm(32 << 20)

	result, err := h.Controller.UpdateRow(tableName, req.Form, id)
	if err != nil {
		log.Error().Err(err).Msg("UpdateEntryHandler")
		http.Error(resp, http.StatusText(500), 500)
		return
	}
	utils.WriteResultToResp(result, resp)
}

func (h *Handlers) DeleteEntryHandler(resp http.ResponseWriter, req *http.Request) {
	log.Info().Msg("Called DeleteEntryHandler")
	tableName := utils.GetVariable("table", req)
	id, err := strconv.Atoi(utils.GetVariable("id", req))
	if err != nil {
		log.Error().Err(err).Msg("")
		http.NotFound(resp, req)
		return
	}

	err = h.Controller.DeleteRow(tableName, id)
	if err != nil {
		log.Error().Err(err).Msg("DeleteEntryHandler")
		http.Error(resp, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(resp, "Row with id = %d was deleted\n", id)
}
