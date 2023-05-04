package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"leafy/app/models"
	"skfw/papaya/pigeon/templates/basicAuth/repository"
)

type UserRepository struct {
	DB *gorm.DB
}

type UserRepositoryImpl interface {
	Init(DB *gorm.DB) error
	SearchFast(username string, email string) (*models.Users, error)
	CatchAll(offset int, limit int) ([]models.Users, error)

	SearchFastById(id uuid.UUID) (*models.Users, error)
	Update(user *models.Users) error

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

		u.DB = DB.Model(&models.Users{})

		return nil
	}

	return errors.New("DB is NULL")
}

func (u *UserRepository) SearchFast(username string, email string) (*models.Users, error) {

	u.NewSession()

	var users []models.Users
	users = make([]models.Users, 0)

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

func (u *UserRepository) SearchFastById(id uuid.UUID) (*models.Users, error) {

	u.NewSession()

	if repository.EmptyIdx(id) {

		return nil, errors.New("id is empty")
	}

	var users []models.Users
	users = make([]models.Users, 0)

	ids := repository.Idx(id)

	if u.DB.Where("id = ?", ids).Limit(1).Find(&users).Error != nil {

		return nil, errors.New("unable to find user")
	}

	if len(users) > 0 {

		return &users[0], nil
	}

	return nil, errors.New("unable to get user information")
}

func (u *UserRepository) CatchAll(offset int, limit int) ([]models.Users, error) {

	u.NewSession()

	var users []models.Users

	if u.DB.Preload("Sessions").Offset(offset).Limit(limit).Find(&users).Error != nil {

		return users, errors.New("unable to find users")
	}

	return users, nil
}

func (u *UserRepository) Update(user *models.Users) error {

	u.NewSession()

	userId := repository.Ids(user.ID)

	if repository.EmptyIdx(userId) {

		return errors.New("invalid id")
	}

	if check, _ := u.SearchFastById(userId); check != nil {

		//check.Name = user.Name
		//check.Gender = user.Gender
		//check.DOB = user.DOB
		//check.Address = user.Address
		//check.Phone = user.Phone
		//check.CountryCode = user.CountryCode
		//check.City = user.City
		//check.PostalCode = user.PostalCode

		if u.DB.Exec("UPDATE users SET name = ?, gender = ?, dob = ?, address = ?, phone = ?, country_code = ?, city = ?, postal_code = ? WHERE id = ?",
			user.Name,
			user.Gender,
			user.DOB,
			user.Address,
			user.Phone,
			user.CountryCode,
			user.City,
			user.PostalCode,
			user.ID,
		).Error != nil {

			return errors.New("unable to update user")
		}

		// why not work as well, idk
		//if err := u.DB.Where("id = ?", userId).Updates(check).Error; err != nil {
		//
		//	return errors.New("unable to update user")
		//}

		return nil
	}

	return errors.New("user not found")
}

func (u *UserRepository) NewSession() {

	u.DB = u.DB.Session(&gorm.Session{})
}
