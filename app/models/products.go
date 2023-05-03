package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	*gorm.Model
	ID          string          `gorm:"type:VARCHAR(32);primary" json:"id"`
	Name        string          `gorm:"not null" json:"name"`
	Image       string          `json:"image"`
	Description string          `json:"description"`
	Stocks      int             `json:"qty"`
	Price       decimal.Decimal `json:"price"`
}

func (Product) TableName() string {

	return "products"
}
