package models

import "gorm.io/gorm"

type Nutrients struct {
	*gorm.Model
	ProductId string `gorm:"type:VARCHAR(32);not null" json:"product_id"`
	Name      string `gorm:"type:VARCHAR(32);not null" json:"name"`
	Unit      string `gorm:"type:VARCHAR(32);not null" json:"unit"`
	Value     int    `gorm:"type:INTEGER;default:0" json:"value"`
}

func (Nutrients) TableName() string {

	return "nutrients"
}
