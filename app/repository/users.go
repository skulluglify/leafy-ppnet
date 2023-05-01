package repository

import (
	"errors"
	"gorm.io/gorm"
	"leafy/app/models"
)

type UserRepository struct {
	DB *gorm.DB
}

type UserRepositoryImpl interface {
	Init(DB *gorm.DB) error
	SearchFast(username string, email string) (*models.User, error)
	NewSession()
}

func UserRepositoryNew(DB *gorm.DB) (UserRepositoryImpl, error) {

	userRepo := &UserRepository{}
	if err := userRepo.Init(DB); err != nil {

		return nil, err
	}
	return userRepo, nil
}

func (u *UserRepository) Init(DB *gorm.DB) error {

	if DB != nil {

		u.DB = DB

		return nil
	}

	return errors.New("DB is NULL")
}

func (u *UserRepository) SearchFast(username string, email string) (*models.User, error) {

	u.NewSession()

	var users []models.User
	users = make([]models.User, 0)

	if username != "" {
		if email != "" {
			if u.DB.Where("username = ? AND email = ?", username, email).Limit(1).Find(&users).Error != nil {

				return nil, errors.New("unable to find user")
			}
		} else {

			if u.DB.Where("username = ?", email).Limit(1).Find(&users).Error != nil {

				return nil, errors.New("unable to find user")
			}
		}
	} else {
		if email != "" {

			if u.DB.Where("email = ?", email).Limit(1).Find(&users).Error != nil {

				return nil, errors.New("unable to find user")
			}
		} else {

			return nil, errors.New("username or email is empty")
		}
	}

	if len(users) > 0 {

		return &users[0], nil
	}

	return nil, errors.New("unable to get user information")
}

func (u *UserRepository) NewSession() {

	u.DB = u.DB.Session(&gorm.Session{})
}
