package repository

import (
	"errors"
	"gorm.io/gorm"
	"leafy/app/models"
)

type AdminRepository struct {
	DB *gorm.DB
}

type AdminRepositoryImpl interface {
	Init(DB *gorm.DB) error
	SearchFast(username string, email string) (*models.User, error)

	CatchUsers(offset int, limit int) ([]models.User, error)

	NewSession()
}

func AdminRepositoryNew(DB *gorm.DB) (AdminRepositoryImpl, error) {

	adminRepo := &AdminRepository{}
	if err := adminRepo.Init(DB); err != nil {

		return nil, err
	}
	return adminRepo, nil
}

func (u *AdminRepository) Init(DB *gorm.DB) error {

	if DB != nil {

		u.DB = DB.Model(&models.User{})

		return nil
	}

	return errors.New("DB is NULL")
}

func (u *AdminRepository) SearchFast(username string, email string) (*models.User, error) {

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

func (u *AdminRepository) CatchUsers(offset int, limit int) ([]models.User, error) {

	u.NewSession()

	var users []models.User

	if u.DB.Preload("Sessions").Offset(offset).Limit(limit).Find(&users).Error != nil {

		return users, errors.New("unable to find users")
	}

	return users, nil
}

func (u *AdminRepository) NewSession() {

	u.DB = u.DB.Session(&gorm.Session{})
}
