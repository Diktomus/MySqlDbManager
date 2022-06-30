package delivery

import (
	"encoding/json"
	"fmt"
	"github/mysql-dbmanager/internal/controller"
	"github/mysql-dbmanager/internal/model"
	"github/mysql-dbmanager/internal/utils"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

type Handlers struct {
	Controller controller.IController
}

func (h *Handlers) GetTablesHandler(resp http.ResponseWriter, req *http.Request) {
	log.Info().Msg("Called GetTablesHandler")
	jsonData, err := json.Marshal(h.Controller.GetTables())
	if err != nil {
		log.Error().Err(err).Msg("GetTablesHandler")
		http.Error(resp, http.StatusText(500), 500)
		return
	}
	fmt.Fprintf(resp, "%s", jsonData)
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

	rows, err := h.Controller.GetRows(tableName, int64(limit), int64(offset))
	if err != nil {
		log.Error().Err(err).Msg("")
		http.NotFound(resp, req)
		return
	}

	jsonData, err := json.Marshal(rows)
	if err != nil {
		log.Error().Err(err).Msg("GetEntriesHandler")
		http.Error(resp, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(resp, "%s", jsonData)
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

	row, err := h.Controller.GetRow(tableName, int64(id))
	if err != nil {
		log.Error().Err(err).Msg("")
		http.NotFound(resp, req)
		return
	}

	jsonData, err := json.Marshal(row)
	if err != nil {
		log.Error().Err(err).Msg("GetEntryHandler")
		http.Error(resp, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(resp, "%s", jsonData)
}

func (h *Handlers) CreateEntryHandler(resp http.ResponseWriter, req *http.Request) {
	log.Info().Msg("Called CreateEntryHandler")
	tableName := utils.GetVariable("table", req)

	req.ParseMultipartForm(32 << 20)
	columns, values, err := utils.ParseUrlValues(req.Form)
	if err != nil {
		log.Error().Err(err).Msg("CreateEntryHandler")
		http.NotFound(resp, req)
		return
	}
	row := model.NewRow(columns, values)
	err = h.Controller.CreateRow(tableName, row)
	if err != nil {
		log.Error().Err(err).Msg("CreateEntryHandler")
		http.Error(resp, http.StatusText(500), 500)
		return
	}
}

func (h *Handlers) UpdateEntryHandler(resp http.ResponseWriter, req *http.Request) {
	log.Info().Msg("Called UpdateEntryHandler")
	tableName := utils.GetVariable("table", req)
	_, err := strconv.Atoi(utils.GetVariable("id", req))
	if err != nil {
		log.Error().Err(err).Msg("")
		http.NotFound(resp, req)
		return
	}

	req.ParseMultipartForm(32 << 20)
	columns, values, err := utils.ParseUrlValues(req.Form)
	if err != nil {
		log.Error().Err(err).Msg("UpdateEntryHandler")
		http.NotFound(resp, req)
		return
	}
	row := model.NewRow(columns, values)
	err = h.Controller.UpdateRow(tableName, row)
	if err != nil {
		log.Error().Err(err).Msg("UpdateEntryHandler")
		http.Error(resp, http.StatusText(500), 500)
		return
	}
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

	err = h.Controller.DeleteRow(tableName, int64(id))
	if err != nil {
		log.Error().Err(err).Msg("DeleteEntryHandler")
		http.Error(resp, http.StatusText(500), 500)
		return
	}

	fmt.Fprintf(resp, "Row with id = %d was deleted\n", id)
}
