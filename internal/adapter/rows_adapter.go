package adapter

import (
	"database/sql"
	"fmt"
	"github/mysql-dbmanager/internal/utils"
	"net/http"
)

type AdaptedRow map[string]interface{}

func AdaptSqlRows(rows *sql.Rows) ([]AdaptedRow, error) {
	adaptedRows := make([]AdaptedRow, 0)
	columns, _ := rows.Columns()
	scanBuffer := utils.ScanBuffer{}
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
	scanBuffer := utils.ScanBuffer{}
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
		} else {
			adaptedRow[column] = values[i]
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
			if columnValue, ok := value.(string); ok {
				columnValues = append(columnValues, columnValue)
			}
		}
	}
	return columnValues
}

func WriteRowsToResp(adaptedRows []AdaptedRow, resp http.ResponseWriter) {
	if len(adaptedRows) > 0 {
		columns := make([]string, 0, len(adaptedRows[0]))
		for column, _ := range adaptedRows[0] {
			columns = append(columns, column)
		}
		fmt.Fprintf(resp, "%+v", column)
	}
	for _, adaptedRow := range adaptedRows {
		for column, value := range adaptedRow {
			fmt.Fprintf(resp, "%+v", adaptedRow)
		}
		fmt.Fprintf(resp, "\n")
	}
}
