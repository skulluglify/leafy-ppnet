package models

import "skfw/papaya/pigeon/templates/basicAuth/models"

type User struct {
	*models.UserModel
	Transactions []Transaction `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"transactions"`
	Carts        []Cart        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"carts"`
}

func (User) TableName() string {

	return "users"
}
