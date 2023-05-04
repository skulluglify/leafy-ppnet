package repository

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"leafy/app/models"
	"skfw/papaya/pigeon/templates/basicAuth/repository"
)

type CartRepository struct {
	DB          *gorm.DB
	productRepo ProductRepositoryImpl
}

type CartRepositoryImpl interface {
	Init(DB *gorm.DB) error

	CatchAll(userId uuid.UUID, offset int, limit int) ([]models.Cart, error)

	SearchFast(userId uuid.UUID, productId uuid.UUID) (*models.Cart, error)
	SearchFastById(id uuid.UUID) (*models.Cart, error)
	SafeSearchFastById(id uuid.UUID, userId uuid.UUID) (*models.Cart, error)

	CreateFast(userId uuid.UUID, productId uuid.UUID, qty int) (*models.Cart, error)

	Update(cart *models.Cart) error
	SafeUpdate(userId uuid.UUID, cart *models.Cart) error

	Delete(id uuid.UUID) error
	SafeDelete(id uuid.UUID, userId uuid.UUID) error

	NewSession()
}

func CartRepositoryNew(DB *gorm.DB) (CartRepositoryImpl, error) {

	cartRepo := &CartRepository{}
	if err := cartRepo.Init(DB); err != nil {

		return nil, err
	}
	return cartRepo, nil
}

func (c *CartRepository) Init(DB *gorm.DB) error {

	var err error

	if DB != nil {

		c.DB = DB.Model(&models.Cart{})
		if c.productRepo, err = ProductRepositoryNew(DB); err != nil {

			return err
		}

		return nil
	}

	return errors.New("DB is NULL")
}

func (c *CartRepository) CatchAll(userId uuid.UUID, offset int, limit int) ([]models.Cart, error) {

	c.NewSession()

	if repository.EmptyIdx(userId) {

		return nil, errors.New("invalid userId")
	}

	userIds := repository.Idx(userId)

	var carts []models.Cart
	carts = make([]models.Cart, 0)

	if c.DB.Where("user_id = ? AND transaction_id IS NULL", userIds).Offset(offset).Limit(limit).Order("created_at DESC").Find(&carts).Error != nil {

		return nil, errors.New("unable to find cart")
	}

	if len(carts) > 0 {

		return carts, nil
	}

	// pass empty array
	//return carts, errors.New("cart is empty")
	return carts, nil
}

func (c *CartRepository) SearchFast(userId uuid.UUID, productId uuid.UUID) (*models.Cart, error) {

	c.NewSession()

	if repository.EmptyIdx(userId) || repository.EmptyIdx(productId) {

		return nil, errors.New("userId or productId is invalid id")
	}

	userIds := repository.Idx(userId)
	productIds := repository.Idx(productId)

	var carts []models.Cart
	carts = make([]models.Cart, 0)

	if c.DB.Where("user_id = ? AND product_id = ? AND transaction_id IS NULL", userIds, productIds).Limit(1).Find(&carts).Error != nil {

		return nil, errors.New("unable to find cart")
	}

	if len(carts) > 0 {

		return &carts[0], nil
	}

	return nil, errors.New("cart is empty")
}

func (c *CartRepository) SearchFastById(id uuid.UUID) (*models.Cart, error) {

	c.NewSession()

	if repository.EmptyIdx(id) {

		return nil, errors.New("invalid id")
	}

	ids := repository.Idx(id)

	var carts []models.Cart
	carts = make([]models.Cart, 0)

	if c.DB.Where("id = ? AND transaction_id IS NULL", ids).Limit(1).Find(&carts).Error != nil {

		return nil, errors.New("unable to find cart")
	}

	if len(carts) > 0 {

		return &carts[0], nil
	}

	return nil, errors.New("cart is empty")
}

func (c *CartRepository) SafeSearchFastById(id uuid.UUID, userId uuid.UUID) (*models.Cart, error) {

	c.NewSession()

	if repository.EmptyIdx(id) || repository.EmptyIdx(userId) {

		return nil, errors.New("id or userId is invalid id")
	}

	ids := repository.Idx(id)
	userIds := repository.Idx(userId)

	var carts []models.Cart
	carts = make([]models.Cart, 0)

	if c.DB.Where("id = ? AND user_id = ? AND transaction_id IS NULL", ids, userIds).Limit(1).Find(&carts).Error != nil {

		return nil, errors.New("unable to find cart")
	}

	if len(carts) > 0 {

		return &carts[0], nil
	}

	return nil, errors.New("cart is empty")
}

func (c *CartRepository) CreateFast(userId uuid.UUID, productId uuid.UUID, qty int) (*models.Cart, error) {

	id := uuid.New()

	if repository.EmptyIdx(userId) || repository.EmptyIdx(productId) {

		return nil, errors.New("userId or productId is invalid id")
	}

	if check, _ := c.SearchFast(userId, productId); check != nil {

		return nil, errors.New("cart has been added")
	}

	var cart models.Cart
	cart = models.Cart{
		ID:            repository.Idx(id),
		UserID:        repository.Idx(userId),
		ProductID:     repository.Idx(productId),
		TransactionID: sql.NullString{Valid: false},
		Qty:           qty,
	}

	// check stocks
	if check, _ := c.productRepo.SearchFastById(productId); check != nil {

		if qty <= check.Stocks {

			if c.DB.Create(&cart).Error != nil {

				return nil, errors.New("unable to create cart")
			}

			return &cart, nil
		}

		return nil, errors.New("quantity out of stocks")
	}

	return nil, errors.New("unable to get product information")
}

func (c *CartRepository) Update(cart *models.Cart) error {

	var err error

	idx := repository.Ids(cart.ID)

	var oldCart *models.Cart

	if oldCart, err = c.SearchFastById(idx); err != nil {

		return err
	}

	// merge
	cart.UserID = oldCart.UserID
	cart.ProductID = oldCart.ProductID

	productId := repository.Ids(cart.ProductID)

	if cart.Qty == 0 {

		cartId := repository.Ids(cart.ID)

		if err = c.Delete(cartId); err != nil {

			return err
		}

		// pass zero qty to delete cart
		return nil
	}

	// check stocks
	if check, _ := c.productRepo.SearchFastById(productId); check != nil {

		if cart.Qty <= check.Stocks {

			if c.DB.Where("id = ?", cart.ID).Updates(cart).Error != nil {

				return errors.New("unable to update cart")
			}

			return nil
		}

		return errors.New("quantity out of stocks")
	}

	return errors.New("unable to get product information")
}

func (c *CartRepository) SafeUpdate(userId uuid.UUID, cart *models.Cart) error {

	var err error

	idx := repository.Ids(cart.ID)
	userIdx := repository.Idx(userId)

	var oldCart *models.Cart

	if oldCart, err = c.SafeSearchFastById(idx, userId); err != nil {

		return err
	}

	// merge
	cart.UserID = oldCart.UserID
	cart.ProductID = oldCart.ProductID

	if cart.Qty == 0 {

		cartId := repository.Ids(cart.ID)

		if err = c.SafeDelete(cartId, userId); err != nil {

			return err
		}

		// pass zero qty to delete cart
		return nil
	}

	productId := repository.Ids(cart.ProductID)

	// check stocks
	if check, _ := c.productRepo.SearchFastById(productId); check != nil {

		if cart.Qty <= check.Stocks {

			if c.DB.Where("id = ? AND user_id = ?", cart.ID, userIdx).Updates(cart).Error != nil {

				return errors.New("unable to update cart")
			}

			return nil
		}

		return errors.New("quantity out of stocks")
	}

	return errors.New("unable to get product information")
}

func (c *CartRepository) Delete(id uuid.UUID) error {

	ids := repository.Idx(id)

	if _, err := c.SearchFastById(id); err != nil {

		return err
	}

	if c.DB.Where("id = ?", ids).Delete(&models.Cart{}).Error != nil {

		return errors.New("unable to delete cart")
	}

	return nil
}

func (c *CartRepository) SafeDelete(id uuid.UUID, userId uuid.UUID) error {

	ids := repository.Idx(id)
	userIds := repository.Idx(userId)

	if _, err := c.SafeSearchFastById(id, userId); err != nil {

		return err
	}

	if c.DB.Where("id = ? AND user_id = ?", ids, userIds).Delete(&models.Cart{}).Error != nil {

		return errors.New("unable to delete cart")
	}

	return nil
}

func (c *CartRepository) NewSession() {

	c.DB = c.DB.Session(&gorm.Session{})
}
