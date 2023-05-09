package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"leafy/app/models"
	bacx "skfw/papaya/pigeon/templates/basicAuth/util"
)

type TransactionRepository struct {
	DB           *gorm.DB
	modelCart    *gorm.DB
	modelProduct *gorm.DB
	cartRepo     CartRepositoryImpl
	productRepo  ProductRepositoryImpl
}

type TransactionRepositoryImpl interface {
	Init(DB *gorm.DB) error

	CatchAll(userId uuid.UUID, offset int, size int) ([]models.Transactions, error)

	SearchFast(cartId uuid.UUID) (*models.Transactions, error)
	SafeSearchFast(userId uuid.UUID, cartId uuid.UUID) (*models.Transactions, error)

	CreateFast(userId uuid.UUID, method string) (*models.Transactions, error)

	Verify(transaction *models.Transactions) error

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
		t.modelCart = DB.Model(&models.Cart{})
		t.modelProduct = DB.Model(&models.Products{})
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

	if bacx.EmptyIdx(userId) {

		return nil, errors.New("invalid userId")
	}

	userIds := bacx.Idx(userId)

	var transactions []models.Transactions

	if t.DB.Unscoped().Preload("Carts").Where("user_id = ?", userIds).Offset(offset).Limit(size).Order("created_at DESC").Find(&transactions).Error != nil {

		return transactions, errors.New("unable to find transactions")
	}

	// pass empty array
	//return transactions, errors.New("transactions is empty")
	return transactions, nil
}

func (t *TransactionRepository) SearchFast(cartId uuid.UUID) (*models.Transactions, error) {

	t.NewSession()

	if bacx.EmptyIdx(cartId) {

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

	if bacx.EmptyIdx(userId) || bacx.EmptyIdx(cartId) {

		return nil, errors.New("userId or cartId is invalid id")
	}

	userIds := bacx.Idx(userId)

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

func (t *TransactionRepository) CreateFast(userId uuid.UUID, method string) (*models.Transactions, error) {

	var err error

	if bacx.EmptyIdx(userId) {

		return nil, errors.New("invalid userId")
	}

	userIds := bacx.Idx(userId)

	// update cart session
	t.modelCart = t.modelCart.Session(&gorm.Session{})

	// the cart is empty

	var count int64

	prepared := t.modelCart.
		Where("carts.user_id = ? AND carts.transaction_id IS NULL", userIds).
		Count(&count)

	if err = prepared.Error; err != nil {

		return nil, errors.New("unable to get cart information")
	}

	if count == 0 {

		return nil, errors.New("cart is empty")
	}

	id := bacx.Idx(uuid.New())

	var transaction models.Transactions
	transaction = models.Transactions{
		ID:            id,
		UserID:        userIds,
		PaymentMethod: method,
	}

	prepared = t.DB.Create(&transaction)

	if err = prepared.Error; err != nil {

		return nil, errors.New("unable to create transaction")
	}

	t.modelProduct = t.modelProduct.Session(&gorm.Session{})

	// update product stocks
	// get all cart id

	type Carts []struct {
		ProductID string
		Qty       int64
	}

	var carts Carts
	carts = make(Carts, 0)

	prepared = t.modelCart.Select("product_id, qty").Where("user_id = ? AND transaction_id IS NULL", userIds).
		Scan(&carts)

	if err = prepared.Error; err != nil {

		return nil, errors.New("unable to catch all carts information")
	}

	for _, cart := range carts {

		// update stocks product
		prepared = t.modelProduct.Exec("UPDATE products SET stocks = stocks - ? WHERE id = ?", cart.Qty, cart.ProductID)
		if err = prepared.Error; err != nil {

			return nil, errors.New("unable to update product stocks")
		}
	}

	// update cart
	prepared = t.modelCart.Where("carts.user_id = ? AND carts.transaction_id IS NULL AND carts.deleted_at IS NULL", userIds).
		Updates(map[string]any{
			"transaction_id": id,
		})

	if err = prepared.Error; err != nil {

		return nil, errors.New("unable to update cart")
	}

	// update cart session
	t.modelCart = t.modelCart.Session(&gorm.Session{})

	// delete cart
	prepared = t.modelCart.Where("carts.user_id = ?", userIds).
		Delete(&models.Cart{})

	if err = prepared.Error; err != nil {

		return nil, errors.New("unable to delete cart")
	}

	return &transaction, err
}

func (t *TransactionRepository) Verify(transaction *models.Transactions) error {

	t.NewSession()

	var err error

	if bacx.EmptyIds(transaction.ID) {

		return errors.New("invalid userId")
	}

	prepared := t.DB.
		Where("id = ? ", transaction.ID).
		Updates(map[string]any{
			"verify": true,
		})

	if err = prepared.Error; err != nil {

		return errors.New("unable to verify transaction, please contact admin")
	}

	return nil
}

func (t *TransactionRepository) NewSession() {

	t.DB = t.DB.Session(&gorm.Session{})
}
