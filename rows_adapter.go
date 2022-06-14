package main

import (
	"database/sql"
)

type AdaptedRow map[string]interface{}

func AdaptSqlRows(rows *sql.Rows) ([]AdaptedRow, error) {
	adaptedRows := make([]AdaptedRow, 0)
	columns, _ := rows.Columns()
	scanBuffer := ScanBuffer{}
	scanBuffer.Prepare(len(columns))

	for rows.Next() {
		err := rows.Scan(scanBuffer.Pointers...)
		if err != nil {
			return nil, err
		}
		adaptedRow := ResolveAdaptedRow(columns, scanBuffer.Values)
		adaptedRows = append(adaptedRows, adaptedRow)
	}
	return adaptedRows, nil
}

func AdaptSqlRow(row *sql.Row, columns []string) (AdaptedRow, error) {
	scanBuffer := ScanBuffer{}
	scanBuffer.Prepare(len(columns))

	err := row.Scan(scanBuffer.Pointers...)
	if err != nil {
		return nil, err
	}

	adaptedRow := ResolveAdaptedRow(columns, scanBuffer.Values)
	return adaptedRow, nil
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

func GetColumnValues(column string, rows *sql.Rows) []string {
	adaptedRows, err := AdaptSqlRows(rows)
	if err != nil {
		return nil
	}

	columnValues := make([]string, 0)
	for _, adaptedRow := range adaptedRows {
		if value, ok := adaptedRow[column]; ok {
			if columnValue, ok := value.([]byte); ok {
				columnValues = append(columnValues, string(columnValue))
			}
		}
	}
	return columnValues
}
