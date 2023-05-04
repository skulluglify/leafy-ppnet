package repository

import (
	"errors"
	"gorm.io/gorm"
	"leafy/app/models"
)

type CategoryRepository struct {
	DB     *gorm.DB
	Shadow *gorm.DB
}

type CategoryRepositoryImpl interface {
	Init(DB *gorm.DB) error
	NewSession()
}

func CategoryRepositoryNew(DB *gorm.DB) (CategoryRepositoryImpl, error) {

	cartRepo := &CategoryRepository{}
	if err := cartRepo.Init(DB); err != nil {

		return nil, err
	}
	return cartRepo, nil
}

func (c *CategoryRepository) Init(DB *gorm.DB) error {

	if DB != nil {

		c.DB = DB.Model(&models.Category{})
		c.Shadow = DB.Model(&models.Categories{})

		return nil
	}

	return errors.New("DB is NULL")
}

func (c *CategoryRepository) NewSession() {

	c.DB = c.DB.Session(&gorm.Session{})
	c.Shadow = c.Shadow.Session(&gorm.Session{})
}
