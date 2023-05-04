package models

import "gorm.io/gorm"

type Categories struct {
	*gorm.Model
	ProductId  string `gorm:"type:VARCHAR(32);not null" json:"product_id"`
	CategoryId string `gorm:"type:VARCHAR(32);not null" json:"category_id"`
}

func (Categories) TableName() string {
	
	return "categories"
}
