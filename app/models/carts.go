package models

import "gorm.io/gorm"

type Cart struct {
	*gorm.Model
	ID            string `gorm:"type:VARCHAR(32);primary" json:"id"`
	UserID        string `gorm:"type:VARCHAR(32);not null" json:"user_id"`
	TransactionID string `gorm:"type:VARCHAR(32)" json:"transaction_id"` // nullable
	ProductID     string `gorm:"type:VARCHAR(32);not null" json:"product_id"`
	Qty           int    `json:"qty"`
}

func (Cart) TableName() string {

	return "carts"
}
