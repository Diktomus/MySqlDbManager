package controller

import "github/mysql-dbmanager/internal/model"

//go:generate mockgen -source=controller.go -destination=mock_controller/mock.go

type IController interface {
	CreateRow(tableName string, newRow model.Row) error
	GetRow(tableName string, rowId int64) (model.Row, error)
	GetRows(tableName string, limit int64, offset int64) ([]model.Row, error)
	UpdateRow(tableName string, newRow model.Row) error
	DeleteRow(tableName string, rowId int64) error
	GetTables() []model.Table
}
