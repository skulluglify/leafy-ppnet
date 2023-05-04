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

	SearchFast(name string) (*models.Products, error)
	CreateFast(name string, description string, price decimal.Decimal, stocks int) (*models.Products, error)

	CatchAll(offset int, limit int) ([]models.Products, error)

	SearchFastById(id uuid.UUID) (*models.Products, error)
	SearchUnscopedFastById(id uuid.UUID) (*models.Products, error)

	Update(id uuid.UUID, product *models.Products) error

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

func (p *ProductRepository) Init(DB *gorm.DB) error {

	if DB != nil {

		p.DB = DB.Model(&models.Products{})

		return nil
	}

	return errors.New("DB is NULL")
}

func (p *ProductRepository) CatchAll(offset int, limit int) ([]models.Products, error) {

	p.NewSession()

	var products []models.Products
	products = make([]models.Products, 0)

	if p.DB.Offset(offset).Limit(limit).Order("created_at DESC").Find(&products).Error != nil {

		return products, errors.New("unable to get product information")
	}

	// pass empty array
	//return products, errors.New("products is empty")
	return products, nil
}

func (p *ProductRepository) SearchFast(name string) (*models.Products, error) {

	p.NewSession()

	if name == "" {

		return nil, errors.New("name is empty")
	}

	var products []models.Products
	products = make([]models.Products, 0)

	if p.DB.Where("name = ?", name).Limit(1).Find(&products).Error != nil {

		return nil, errors.New("unable to search product")
	}

	if len(products) > 0 {

		return &products[0], nil
	}

	return nil, errors.New("product is empty")
}

func (p *ProductRepository) SearchFastById(id uuid.UUID) (*models.Products, error) {

	p.NewSession()

	if repository.EmptyIdx(id) {

		return nil, errors.New("id is empty")
	}

	var products []models.Products
	products = make([]models.Products, 0)

	if p.DB.Where("id = ?", repository.Idx(id)).Limit(1).Find(&products).Error != nil {

		return nil, errors.New("unable to search product")
	}

	if len(products) > 0 {

		return &products[0], nil
	}

	return nil, nil
}

func (p *ProductRepository) SearchUnscopedFastById(id uuid.UUID) (*models.Products, error) {

	p.NewSession()

	if repository.EmptyIdx(id) {

		return nil, errors.New("id is empty")
	}

	var products []models.Products
	products = make([]models.Products, 0)

	// find all products
	if p.DB.Unscoped().Where("id = ?", repository.Idx(id)).Limit(1).Find(&products).Error != nil {

		return nil, errors.New("unable to search product")
	}

	if len(products) > 0 {

		return &products[0], nil
	}

	return nil, nil
}

func (p *ProductRepository) CreateFast(name string, description string, price decimal.Decimal, stocks int) (*models.Products, error) {

	if check, _ := p.SearchFast(name); check != nil {

		return nil, errors.New("product has been added")
	}

	product := &models.Products{
		ID:          repository.Idx(uuid.New()),
		Name:        name,
		Description: description,
		Price:       price,
		Stocks:      stocks,
	}

	if p.DB.Create(&product).Error != nil {

		return product, errors.New("unable to create product")
	}

	return product, nil
}

func (p *ProductRepository) Update(id uuid.UUID, product *models.Products) error {

	if repository.EmptyIdx(id) {

		return errors.New("id is empty")
	}

	// merge Id
	product.ID = repository.Idx(id)

	if check, _ := p.SearchFastById(id); check != nil {

		if p.DB.Where("id = ?", product.ID).Updates(product).Error != nil {

			return errors.New("unable to update product")
		}
	}

	return nil
}

func (p *ProductRepository) DeleteFast(id uuid.UUID) error {

	if repository.EmptyIdx(id) {

		return errors.New("id is empty")
	}

	// merge Id
	ids := repository.Idx(id)

	if check, _ := p.SearchFastById(id); check != nil {

		if p.DB.Where("id = ?", ids).Delete(&models.Products{}).Error != nil {

			return errors.New("unable to update product")
		}
	}

	return nil
}

func (p *ProductRepository) NewSession() {

	p.DB = p.DB.Session(&gorm.Session{})
}
