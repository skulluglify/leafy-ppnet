package repository

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"leafy/app/models"
	"skfw/papaya/pigeon/templates/basicAuth/repository"
)

type TransactionRepository struct {
	DB          *gorm.DB
	cartRepo    CartRepositoryImpl
	productRepo ProductRepositoryImpl
}

type TransactionRepositoryImpl interface {
	Init(DB *gorm.DB) error

	CatchAll(userId uuid.UUID, offset int, size int) ([]models.Transactions, error)

	SearchFast(cartId uuid.UUID) (*models.Transactions, error)
	SafeSearchFast(userId uuid.UUID, cartId uuid.UUID) (*models.Transactions, error)

	CreateFast(userId uuid.UUID, cartId uuid.UUID, method string) (*models.Transactions, error)

	NewSession()
}

func TransactionRepositoryNew(DB *gorm.DB) (TransactionRepositoryImpl, error) {

	transactionRepo := &TransactionRepository{}
	if err := transactionRepo.Init(DB); err != nil {

		return nil, err
	}
	return transactionRepo, nil
}

func (t *TransactionRepository) Init(DB *gorm.DB) error {

	var err error

	if DB != nil {

		t.DB = DB.Model(&models.Transactions{})
		if t.cartRepo, err = CartRepositoryNew(DB); err != nil {

			return err
		}
		if t.productRepo, err = ProductRepositoryNew(DB); err != nil {

			return err
		}

		return nil
	}

	return errors.New("DB is NULL")
}

func (t *TransactionRepository) CatchAll(userId uuid.UUID, offset int, size int) ([]models.Transactions, error) {

	t.NewSession()

	if repository.EmptyIdx(userId) {

		return nil, errors.New("invalid userId")
	}

	userIds := repository.Idx(userId)

	var transactions []models.Transactions

	if t.DB.Preload("Carts").Where("user_id = ?", userIds).Offset(offset).Limit(size).Order("created_at DESC").Find(&transactions).Error != nil {

		return transactions, errors.New("unable to find transactions")
	}

	// pass empty array
	//return transactions, errors.New("transactions is empty")
	return transactions, nil
}

func (t *TransactionRepository) SearchFast(cartId uuid.UUID) (*models.Transactions, error) {

	t.NewSession()

	if repository.EmptyIdx(cartId) {

		return nil, errors.New("invalid cartId")
	}

	var transactions []models.Transactions
	transactions = make([]models.Transactions, 0)

	var err error
	var cart *models.Cart

	if cart, err = t.cartRepo.SearchFastById(cartId); cart != nil {

		ids := cart.TransactionID

		if ids.Valid {

			if t.DB.Where("id = ?", ids.String).Limit(1).Find(&transactions).Error != nil {

				return nil, errors.New("unable to find transaction")
			}

			if len(transactions) > 0 {

				return &transactions[0], nil
			}

			return nil, errors.New("transaction is empty")
		}

		return nil, errors.New("transaction is not created")
	}

	return nil, err
}

func (t *TransactionRepository) SafeSearchFast(userId uuid.UUID, cartId uuid.UUID) (*models.Transactions, error) {

	t.NewSession()

	if repository.EmptyIdx(userId) || repository.EmptyIdx(cartId) {

		return nil, errors.New("userId or cartId is invalid id")
	}

	userIds := repository.Idx(userId)

	var transactions []models.Transactions
	transactions = make([]models.Transactions, 0)

	var err error
	var cart *models.Cart

	if cart, err = t.cartRepo.SearchFastById(cartId); cart != nil {

		ids := cart.TransactionID

		if ids.Valid {

			if t.DB.Where("id = ? AND user_id = ?", ids.String, userIds).Limit(1).Find(&transactions).Error != nil {

				return nil, errors.New("unable to find transaction")
			}

			if len(transactions) > 0 {

				return &transactions[0], nil
			}

			return nil, errors.New("transaction is empty")

		}

		return nil, errors.New("transaction is not created")
	}

	return nil, err
}

func (t *TransactionRepository) CreateFast(userId uuid.UUID, cartId uuid.UUID, method string) (*models.Transactions, error) {

	var err error

	if repository.EmptyIdx(userId) || repository.EmptyIdx(cartId) {

		return nil, errors.New("userId or cartId is invalid id")
	}

	var check *models.Transactions
	if check, err = t.SafeSearchFast(userId, cartId); check != nil {

		// if cart not invisible, but registered on transaction
		return nil, errors.New("transaction has been added")
	}

	// transaction makes cart invisible to find it > cart is empty
	if err.Error() == "cart is empty" {

		return nil, err
	}

	id := repository.Idx(uuid.New())

	userIds := repository.Idx(userId)

	var transaction models.Transactions
	transaction = models.Transactions{
		ID:            id,
		UserID:        userIds,
		PaymentMethod: method,
	}

	var cart *models.Cart

	if cart, err = t.cartRepo.SafeSearchFastById(cartId, userId); cart != nil {

		productId := repository.Ids(cart.ProductID)

		// search product
		if product, _ := t.productRepo.SearchFastById(productId); product != nil {

			if product.Stocks < cart.Qty {

				return nil, errors.New("quantity out of stocks")
			}

			// desc product stocks
			product.Stocks += -cart.Qty

			// update product
			if err = t.productRepo.Update(productId, product); err != nil {

				return nil, err
			}

			// create transaction
			if t.DB.Create(&transaction).Error != nil {

				return nil, errors.New("unable to create transaction")
			}

			// merge
			cart.TransactionID = sql.NullString{String: id, Valid: true}

			// update cart
			if err = t.cartRepo.Update(cart); err != nil {

				return nil, err
			}

			return &transaction, nil
		}

		return nil, errors.New("unable to get product information")
	}

	return nil, err
}

func (t *TransactionRepository) NewSession() {

	t.DB = t.DB.Session(&gorm.Session{})
}
