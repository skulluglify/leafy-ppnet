package repository

import (
	"errors"
	"gorm.io/gorm"
	"leafy/app/models"
)

type NutrientRepository struct {
	DB *gorm.DB
}

type NutrientRepositoryImpl interface {
	Init(DB *gorm.DB) error
	NewSession()
}

func NutrientRepositoryNew(DB *gorm.DB) (NutrientRepositoryImpl, error) {

	nutrientRepo := &NutrientRepository{}
	if err := nutrientRepo.Init(DB); err != nil {

		return nil, err
	}
	return nutrientRepo, nil
}

func (n *NutrientRepository) Init(DB *gorm.DB) error {

	if DB != nil {

		n.DB = DB.Model(&models.Nutrients{})

		return nil
	}

	return errors.New("DB is NULL")
}

func (n *NutrientRepository) NewSession() {

	n.DB = n.DB.Session(&gorm.Session{})
}
