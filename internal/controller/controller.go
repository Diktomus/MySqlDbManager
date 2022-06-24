package controller

import "github/mysql-dbmanager/internal/model"

type RowsController interface {
	CreateRow(tableName string, newRow model.Row) error
	GetRow(tableName string, rowId int64) (model.Row, error)
	GetRows(tableName string, limit int64, offset int64) (model.Row, error)
	UpdateRow(tableName string, newRow model.Row) error
	DeleteRow(tableName string, rowId int64) error
}
