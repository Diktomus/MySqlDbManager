package controller

import (
	"database/sql"
	"fmt"
	"github/mysql-dbmanager/internal/model"
	"github/mysql-dbmanager/internal/utils"
	"strings"

	"github.com/rs/zerolog/log"
)

type MySqlRowsController struct {
	db     *sql.DB
	tables []model.Table
}

func (controller *MySqlRowsController) getColumns(tableName string) []string {
	for _, table := range controller.tables {
		if table.Name == tableName {
			return table.Columns
		}
	}
	return []string{}
}

func (controller *MySqlRowsController) GetRows(tableName string, limit int64, offset int64) ([]model.Row, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id >= ? LIMIT ?", tableName)
	sqlRows, err := controller.db.Query(query, offset, limit)
	defer sqlRows.Close()
	if err != nil {
		return nil, err
	}
	rows := make([]model.Row, 0)
	columns, _ := sqlRows.Columns()
	scanBuffer := utils.ScanBuffer{}
	scanBuffer.Prepare(len(columns))

	rowId := offset
	for sqlRows.Next() {
		err := sqlRows.Scan(scanBuffer.Pointers...)
		if err != nil {
			return nil, err
		}
		row := model.NewRow(rowId, columns, scanBuffer.Values)
		rows = append(rows, row)
		rowId++
	}

	return rows, nil
}

func (controller *MySqlRowsController) GetRow(tableName string, rowId int64) (model.Row, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", tableName)
	sqlRow := controller.db.QueryRow(query, rowId)

	scanBuffer := utils.ScanBuffer{}
	columns := controller.getColumns(tableName)
	scanBuffer.Prepare(len(columns))

	err := sqlRow.Scan(scanBuffer.Pointers...)
	if err != nil {
		return model.Row{}, err
	}

	row := model.NewRow(rowId, columns, scanBuffer.Values)
	return row, nil
}

func (controller *MySqlRowsController) UpdateRow(tableName string, newRow model.Row) error {
	columns, values := controller.getColumnsAndValues(tableName, newRow)
	values = append(values, newRow.Id)
	columns = append(columns, "")
	queryColumnsPlaceholders := strings.Join(columns, " = ?, ")
	queryColumnsPlaceholders = strings.TrimRight(queryColumnsPlaceholders, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", tableName, queryColumnsPlaceholders)
	_, err := controller.db.Exec(query, values...)
	if err != nil {
		return err
	}

	return nil
}

func (controller *MySqlRowsController) CreateRow(tableName string, newRow model.Row) error {
	columns, values := controller.getColumnsAndValues(tableName, newRow)
	queryColumns := strings.Join(columns, ", ")
	queryPlaceholders := strings.Repeat("?, ", len(columns))
	queryPlaceholders = strings.TrimRight(queryPlaceholders, ", ")

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", tableName, queryColumns, queryPlaceholders)
	_, err := controller.db.Exec(query, values...)
	if err != nil {
		return err
	}

	return err
}

func (controller *MySqlRowsController) DeleteRow(tableName string, rowId int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)
	_, err := controller.db.Exec(query, rowId)
	if err != nil {
		return err
	}
	return nil
}

func (controller *MySqlRowsController) GetTables() []model.Table {
	return controller.tables
}

func (controller *MySqlRowsController) getColumnsAndValues(tableName string, row model.Row) (columns []string, values []interface{}) {
	values = make([]interface{}, 0)
	columns = make([]string, 0)
	for _, column := range controller.getColumns(tableName) {
		if value, ok := row.ColumnsByValues[column]; ok {
			values = append(values, value)
			columns = append(columns, column)
		}
	}
	return columns, values
}

func (controller *MySqlRowsController) getFirstColumnValues(query string) ([]string, error) {
	sqlRows, err := controller.db.Query(query)
	defer sqlRows.Close()
	if err != nil {
		return nil, err
	}
	columns, _ := sqlRows.Columns()
	scanBuffer := utils.ScanBuffer{}
	scanBuffer.Prepare(len(columns))
	firstColumnValues := make([]string, 0)
	for sqlRows.Next() {
		err := sqlRows.Scan(scanBuffer.Pointers...)
		if err != nil {
			return nil, err
		}
		if value, ok := scanBuffer.Values[0].([]byte); ok {
			firstColumnValues = append(firstColumnValues, string(value))
		}
	}

	return firstColumnValues, nil
}

func (controller *MySqlRowsController) Init() error {
	tablesNames, err := controller.getFirstColumnValues("SHOW TABLES")
	if err != nil {
		log.Error().Err(err).Msg("controller.Init.getTablesNames:")
		return err
	}
	tables := make([]model.Table, 0, len(tablesNames))

	for _, tableName := range tablesNames {
		query := fmt.Sprintf("SHOW COLUMNS FROM %s", tableName)
		columns, err := controller.getFirstColumnValues(query)
		if err != nil {
			log.Error().Err(err).Msg("controller.Init.getFirstColumnValues:")
			return err
		}
		table := model.Table{Name: tableName, Columns: columns}
		tables = append(tables, table)
	}
	controller.tables = tables

	return nil
}

func NewMySqlRowsController(database *sql.DB) *MySqlRowsController {
	return &MySqlRowsController{db: database}
}
