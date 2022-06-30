package model

type Row struct {
	ColumnsByValues map[string]interface{} `json:"columnsByValues"`
}

func (row *Row) GetId() int64 {
	if id, ok := row.ColumnsByValues["id"].(int64); ok {
		return id
	}
	return 0
}

func NewRow(columns []string, values []interface{}) Row {
	row := &Row{
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
