package models

import (
	"github.com/shopspring/decimal"
	"skfw/papaya/pigeon/templates/basicAuth/models"
)

type Users struct {
	*models.UserModel
	Balance      decimal.Decimal `gorm:"default:0" json:"balance"`
	Transactions []Transactions  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"transactions"`
	Carts        []Cart          `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"carts"`
}

func (Users) TableName() string {

	return "users"
}
