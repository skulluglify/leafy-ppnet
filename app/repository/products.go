package repository

import (
	"errors"
	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

type ProductRepositoryImpl interface {
	Init(DB *gorm.DB) error
	NewSession()
}

func ProductRepositoryNew(DB *gorm.DB) (ProductRepositoryImpl, error) {

	productRepo := &ProductRepository{}
	if err := productRepo.Init(DB); err != nil {

		return nil, err
	}
	return productRepo, nil
}

func (u *ProductRepository) Init(DB *gorm.DB) error {

	if DB != nil {

		u.DB = DB

		return nil
	}

	return errors.New("DB is NULL")
}

func (u *ProductRepository) NewSession() {

	u.DB = u.DB.Session(&gorm.Session{})
}
