package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"leafy/app/models"
	"skfw/papaya/pigeon/templates/basicAuth/repository"
)

type ProductRepository struct {
	DB *gorm.DB
}

type ProductRepositoryImpl interface {
	Init(DB *gorm.DB) error

	SearchFast(name string) (*models.Product, error)
	CreateFast(name string, description string, price decimal.Decimal, stocks int) (*models.Product, error)

	CatchAllProducts(offset int, limit int) ([]models.Product, error)

	SearchFastById(id uuid.UUID) (*models.Product, error)

	Update(id uuid.UUID, product *models.Product) error

	DeleteFast(id uuid.UUID) error

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

		u.DB = DB.Model(&models.Product{})

		return nil
	}

	return errors.New("DB is NULL")
}

func (u *ProductRepository) CatchAllProducts(offset int, limit int) ([]models.Product, error) {

	var products []models.Product
	products = make([]models.Product, 0)

	if u.DB.Offset(offset).Limit(limit).Find(&products).Error != nil {

		return products, errors.New("unable to get product information")
	}

	return products, nil
}

func (u *ProductRepository) SearchFast(name string) (*models.Product, error) {

	u.NewSession()

	if name == "" {

		return nil, errors.New("name is empty")
	}

	var products []models.Product
	products = make([]models.Product, 0)

	if u.DB.Where("name = ?", name).Limit(1).Find(&products).Error != nil {

		return nil, errors.New("unable to search product")
	}

	if len(products) > 0 {

		return &products[0], nil
	}

	return nil, nil
}

func (u *ProductRepository) SearchFastById(id uuid.UUID) (*models.Product, error) {

	u.NewSession()

	if repository.EmptyIdx(id) {

		return nil, errors.New("id is empty")
	}

	var products []models.Product
	products = make([]models.Product, 0)

	if u.DB.Where("id = ?", repository.Idx(id)).Limit(1).Find(&products).Error != nil {

		return nil, errors.New("unable to search product")
	}

	if len(products) > 0 {

		return &products[0], nil
	}

	return nil, nil
}

func (u *ProductRepository) CreateFast(name string, description string, price decimal.Decimal, stocks int) (*models.Product, error) {

	if check, _ := u.SearchFast(name); check != nil {

		return nil, errors.New("product has been added")
	}

	product := &models.Product{
		ID:          repository.Idx(uuid.New()),
		Name:        name,
		Description: description,
		Price:       price,
		Stocks:      stocks,
	}

	if u.DB.Create(&product).Error != nil {

		return product, errors.New("unable to create product")
	}

	return product, nil
}

func (u *ProductRepository) Update(id uuid.UUID, product *models.Product) error {

	if repository.EmptyIdx(id) {

		return errors.New("id is empty")
	}

	// merge Id
	product.ID = repository.Idx(id)

	if check, _ := u.SearchFastById(id); check != nil {

		if u.DB.Where("id = ?", product.ID).Updates(product).Error != nil {

			return errors.New("unable to update product")
		}
	}

	return nil
}

func (u *ProductRepository) DeleteFast(id uuid.UUID) error {

	if repository.EmptyIdx(id) {

		return errors.New("id is empty")
	}

	// merge Id
	ids := repository.Idx(id)

	if check, _ := u.SearchFastById(id); check != nil {

		if u.DB.Where("id = ?", ids).Delete(&models.Product{}).Error != nil {

			return errors.New("unable to update product")
		}
	}

	return nil
}

func (u *ProductRepository) NewSession() {

	u.DB = u.DB.Session(&gorm.Session{})
}
