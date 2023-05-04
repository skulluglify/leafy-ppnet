package models

import "skfw/papaya/pigeon/templates/basicAuth/models"

type Users struct {
	*models.UserModel
	Transactions []Transactions `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"transactions"`
	Carts        []Cart         `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"carts"`
}

func (Users) TableName() string {

	return "users"
}
