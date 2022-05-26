package models

import (
	"database/sql"
	"time"
)

// type for database connection values
type DBModel struct {
	DB *sql.DB
}

// wrapper for all odels
type Models struct {
	DB DBModel
}

// returns a model type with db connection pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

// type for product information consistency
type Widget struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}
