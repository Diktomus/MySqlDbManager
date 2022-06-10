package main

import (
	"database/sql"
	"fmt"
	"strings"
)

type Table struct {
	Name    string
	Columns []string
}

type Controller struct {
	DB     *sql.DB
	Tables []Table
}

func (controller *Controller) GetColumns(tableName string) []string {
	for _, table := range controller.Tables {
		if table.Name == tableName {
			return table.Columns
		}
	}
	return []string{}
}

func (controller *Controller) GetRows(tableName string, limit int, offset int) *sql.Rows {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id >= ? LIMIT ?", tableName)
	rows, err := controller.DB.Query(query, offset, limit)
	if err != nil {
		panic(err)
	}
	return rows
}

func (controller *Controller) GetRow(tableName string, id int) *sql.Row {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", tableName)
	row := controller.DB.QueryRow(query, id)
	return row
}

func (controller *Controller) UpdateRow(tableName string, columnsByValues map[string][]string, id int) sql.Result {
	queryArgs, updatingColumns := controller.resolveQueryParams(tableName, columnsByValues)
	queryArgs = append(queryArgs, id)
	updatingColumns = append(updatingColumns, "")
	queryArgsPlaceholders := strings.Join(updatingColumns, " = ?, ")
	queryArgsPlaceholders = strings.TrimRight(queryArgsPlaceholders, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", tableName, queryArgsPlaceholders)
	result, err := controller.DB.Exec(query, queryArgs...)
	if err != nil {
		panic(err)
	}

	return result
}

func (controller *Controller) CreateRow(tableName string, columnsByValues map[string][]string) sql.Result {
	queryArgs, columns := controller.resolveQueryParams(tableName, columnsByValues)
	queryArgsNames := strings.Join(columns, ", ")
	queryPlaceholders := strings.Repeat("?, ", len(columns))
	queryPlaceholders = strings.TrimRight(queryPlaceholders, ", ")

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", tableName, queryArgsNames, queryPlaceholders)
	result, err := controller.DB.Exec(query, queryArgs...)
	if err != nil {
		panic(err)
	}

	return result
}

func (controller *Controller) DeleteRow(tableName string, id int) {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ?", tableName)
	_, err := controller.DB.Exec(query, id)
	if err != nil {
		panic(err)
	}
}

func (controller *Controller) resolveQueryParams(tableName string, colsByVals map[string][]string) ([]interface{}, []string) {
	queryArgs := make([]interface{}, 0)
	affectedColumns := make([]string, 0)
	for _, column := range controller.GetColumns(tableName) {
		if values, ok := colsByVals[column]; ok {
			if len(values) > 0 {
				queryArg := ParseStr(values[0])
				queryArgs = append(queryArgs, queryArg)
				affectedColumns = append(affectedColumns, column)
			}
		}
	}
	return queryArgs, affectedColumns
}

func NewController(database *sql.DB) *Controller {
	rows, err := database.Query("SHOW TABLES")
	if err != nil {
		panic(err)
	}
	tablesNames := make([]string, 0)
	var value string
	for rows.Next() {
		err = rows.Scan(&value)
		if err != nil {
			panic(err)
		}
		tablesNames = append(tablesNames, value)
	}

	tables := make([]Table, 0, len(tablesNames))

	for _, tableName := range tablesNames {
		query := fmt.Sprintf("SHOW COLUMNS FROM %s", tableName)
		rows, err = database.Query(query)
		if err != nil {
			panic(err)
		}

		queryColumns, _ := rows.Columns()
		scanBuffer := PrepareScanBuffer(len(queryColumns))
		columns := make([]string, 0)
		for rows.Next() {
			err = rows.Scan(scanBuffer.Pointers...)
			if err != nil {
				panic(err)
			}

			if column, ok := scanBuffer.Values[0].([]byte); ok {
				columns = append(columns, string(column))
			}
		}
		table := Table{Name: tableName, Columns: columns}
		tables = append(tables, table)
	}

	rows.Close()

	return &Controller{DB: database, Tables: tables}
}
