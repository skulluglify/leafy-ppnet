package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

	PayCheck(id uuid.UUID) error

	Balance(id uuid.UUID) (decimal.Decimal, error)
	Bill(id uuid.UUID) (decimal.Decimal, error)

	Topup(id uuid.UUID, balance decimal.Decimal) error

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

		prepared := u.DB.Where("id = ?", user.ID).Updates(map[string]any{
			"name":         user.Name,
			"gender":       user.Gender,
			"dob":          user.DOB,
			"address":      user.Address,
			"phone":        user.Phone,
			"country_code": user.CountryCode,
			"city":         user.City,
			"postal_code":  user.PostalCode,
		})

		//fmt.Println(prepared.ToSQL(func(tx *gorm.DB) *gorm.DB {
		//	return tx.Model(&models.Users{}).Where("id = ?", user.ID).Updates(map[string]any{
		//		"name":         user.Name,
		//		"gender":       user.Gender,
		//		"dob":          user.DOB,
		//		"address":      user.Address,
		//		"phone":        user.Phone,
		//		"country_code": user.CountryCode,
		//		"city":         user.City,
		//		"postal_code":  user.PostalCode,
		//	})
		//}))

		if err := prepared.Error; err != nil {

			return errors.New("unable to update user")
		}

		//fmt.Println(prepared.RowsAffected)

		//var count int64
		//u.DB.Where("id = ?", userId).Count(&count)
		//fmt.Printf("%d rows matched\n", count)

		return nil
	}

	return errors.New("user not found")
}

func (u *UserRepository) PayCheck(id uuid.UUID) error {

	u.NewSession()

	var err error

	if repository.EmptyIdx(id) {

		return errors.New("invalid id")
	}

	ids := repository.Idx(id)

	var balance decimal.Decimal

	if balance, err = u.Balance(id); err != nil {

		return err
	}

	// get bill
	var bill decimal.Decimal
	if bill, err = u.Bill(id); err != nil {

		return errors.New("unable to get user bill information")
	}

	// check balance
	if bill.GreaterThan(balance) {

		return errors.New("user balance not enough to pay bill")
	}

	// pay
	prepared := u.DB.
		Where("users.id = ?", ids).
		Updates(map[string]any{
			"balance": balance.Sub(bill),
		})

	if err = prepared.Error; err != nil {

		return errors.New("unable to pay bill")
	}

	return nil
}

func (u *UserRepository) Balance(id uuid.UUID) (decimal.Decimal, error) {

	u.NewSession()

	noop := decimal.New(0, 0)

	if repository.EmptyIdx(id) {

		return noop, errors.New("invalid id")
	}

	var err error

	ids := repository.Idx(id)

	type Row struct {
		Balance decimal.Decimal
	}

	var row Row
	row = Row{}

	prepared := u.DB.Where("id = ?", ids).
		Limit(1).
		Scan(&row)

	if err = prepared.Error; err != nil {

		return noop, errors.New("unable to get user balance information")
	}

	return row.Balance, nil
}

func (u *UserRepository) Bill(id uuid.UUID) (decimal.Decimal, error) {

	u.NewSession()

	noop := decimal.New(0, 0)

	if repository.EmptyIdx(id) {

		return noop, errors.New("invalid id")
	}

	var err error

	ids := repository.Idx(id)
	bill := decimal.New(0, 0)

	type Rows []struct {
		Qty   int64
		Price decimal.Decimal
	}

	var rows Rows
	rows = make(Rows, 0)

	prepared := u.DB.
		Select("carts.qty AS qty, products.price AS price").
		Joins("INNER JOIN carts ON users.id = carts.user_id").
		Joins("INNER JOIN products ON carts.product_id = products.id").
		Where("users.id = ? AND carts.transaction_id IS NULL AND carts.deleted_at IS NULL", ids).
		Scan(&rows)

	if err = prepared.Error; err != nil {

		return noop, errors.New("unable to get user information")
	}

	for _, row := range rows {

		var q decimal.Decimal
		q = decimal.NewFromInt(row.Qty)
		bill = bill.Add(q.Mul(row.Price))
	}

	//db, _ := u.DB.DB()
	//
	//var rows *sql.Rows
	//rows, err = db.Query("SELECT carts.qty AS qty, products.price AS price FROM users INNER JOIN carts ON users.id = carts.user_id INNER JOIN products ON carts.product_id = products.id WHERE users.id = $1", ids)
	//
	//if err != nil {
	//
	//	return noop, errors.New("unable to get value")
	//}
	//
	//for rows.Next() {
	//
	//	var qty int64
	//	var price string
	//
	//	if err = rows.Scan(&qty, &price); err != nil {
	//
	//		return noop, errors.New("unable to parse value")
	//	}
	//
	//	var p decimal.Decimal
	//
	//	quantity := decimal.NewFromInt(qty)
	//	if p, err = decimal.NewFromString(price); err != nil {
	//
	//		return noop, errors.New("unable to parse price value")
	//	}
	//
	//	bill = bill.Add(quantity.Mul(p))
	//}
	//
	//if err = rows.Err(); err != nil {
	//
	//	return noop, errors.New("error during parse value in iteration")
	//}

	return bill, nil
}

func (u *UserRepository) Topup(id uuid.UUID, balance decimal.Decimal) error {

	var err error

	u.NewSession()

	if repository.EmptyIdx(id) {

		return errors.New("invalid userId")
	}

	ids := repository.Idx(id)

	var currBalance decimal.Decimal

	if currBalance, err = u.Balance(id); err != nil {

		return err
	}

	prepared := u.DB.Where("id = ?", ids).Updates(map[string]any{
		"balance": currBalance.Add(balance),
	})

	if err = prepared.Error; err != nil {

		return errors.New("unable to update user balance")
	}

	return nil
}

func (u *UserRepository) NewSession() {

	u.DB = u.DB.Session(&gorm.Session{})
}
