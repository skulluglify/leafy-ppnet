package models

import "gorm.io/gorm"

type Category struct {
	*gorm.Model
	ID          string `gorm:"type:VARCHAR(32);primary" json:"id"`
	Name        string `gorm:"type:VARCHAR(32);not null" json:"name"`
	Description string `json:"description"`
}

func (Category) TableName() string {

	return "category"
}
