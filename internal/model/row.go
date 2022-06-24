package model

type Row struct {
	Id              int64
	ColumnsByValues map[string]interface{}
}

func NewRow(rowId int64, columns []string, values []interface{}) Row {
	row := &Row{
		Id:              rowId,
		ColumnsByValues: make(map[string]interface{}, 0),
	}
	for i, column := range columns {
		if value, ok := (values[i]).([]byte); ok {
			row.ColumnsByValues[column] = string(value)
		} else {
			row.ColumnsByValues[column] = values[i]
		}
	}
	return *row
}
