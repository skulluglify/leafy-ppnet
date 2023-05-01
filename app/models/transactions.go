package models

import "gorm.io/gorm"

type Transaction struct {
	*gorm.Model
	ID            string `gorm:"type:VARCHAR(32);primary" json:"id"`
	UserID        string `gorm:"type:VARCHAR(32);not null" json:"user_id"`
	Carts         []Cart `gorm:"foreignKey:TransactionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"carts"`
	PaymentMethod string `json:"payment_method"`
	Verify        bool   `json:"verify"`
}

func (Transaction) TableName() string {

	return "transactions"
}
