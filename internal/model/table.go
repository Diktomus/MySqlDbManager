package model

type Table struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
}
