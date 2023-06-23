package entity

import (
	"context"
	"database/sql"
)

type OrderRepositoryInterface interface {
	Save(Order *Order) error
	List() ([]Order)
	// GetTotal() (int, error)
}

type DBTX interface {
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}
