package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"leafy/app/models"
	m "skfw/papaya/koala/mapping"
	"skfw/papaya/koala/pp"
	bacx "skfw/papaya/pigeon/templates/basicAuth/util"
)

type CategoryRepository struct {
	DB     *gorm.DB
	Shadow *gorm.DB
}

type CategoryRepositoryImpl interface {
	Init(DB *gorm.DB) error
	SearchFast(name string) (*models.Category, error)
	SearchFastById(id uuid.UUID) (*models.Category, error)
	CreateFast(name string, description string) (*models.Category, error)
	Add(productId uuid.UUID, name string) error
	Categories(productId uuid.UUID) []string
	CatchAll(products []models.Products) []m.KMapImpl
	UpdateByName(name string, newName string) error
	DeleteByName(name string) error
	UnlinkByProductId(productId uuid.UUID) error
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

func (c *CategoryRepository) SearchFast(name string) (*models.Category, error) {

	var err error

	c.NewSession()

	if name == "" {

		return nil, errors.New("name is empty string")
	}

	var categories []models.Category

	if err = c.DB.Where("name = ?", name).Limit(1).Find(&categories).Error; err != nil {

		return nil, errors.New("unable to find category")
	}

	if len(categories) > 0 {

		return &categories[0], nil
	}

	return nil, errors.New("category not found")
}

func (c *CategoryRepository) SearchFastById(id uuid.UUID) (*models.Category, error) {

	var err error

	c.NewSession()

	if bacx.EmptyIdx(id) {

		return nil, errors.New("invalid id")
	}

	ids := bacx.Idx(id)

	var categories []models.Category

	if err = c.DB.Where("id = ?", ids).Limit(1).Find(&categories).Error; err != nil {

		return nil, errors.New("unable to find category")
	}

	if len(categories) > 0 {

		return &categories[0], nil
	}

	return nil, errors.New("category not found")
}

func (c *CategoryRepository) CreateFast(name string, description string) (*models.Category, error) {

	var err error

	if name == "" {

		return nil, errors.New("name is empty string")
	}

	if check, _ := c.SearchFast(name); check != nil {

		return nil, errors.New("category has been added")
	}

	category := models.Category{
		ID:          bacx.Idx(uuid.New()),
		Name:        name,
		Description: description,
	}

	if err = c.DB.Create(&category).Error; err != nil {

		return nil, errors.New("unable to create category")
	}

	return &category, nil
}

func (c *CategoryRepository) Add(productId uuid.UUID, name string) error {

	var err error

	if bacx.EmptyIdx(productId) {

		return errors.New("invalid productId")
	}

	var check *models.Category

	productIds := bacx.Idx(productId)

	categories := models.Categories{
		ProductId: productIds,
	}

	// try find - linked
	if check, _ = c.SearchFast(name); check != nil {

		categories.CategoryId = check.ID

		// try - find
		var cats []models.Categories
		cats = make([]models.Categories, 0)

		if err = c.Shadow.Where("category_id = ? AND product_id = ?", check.ID, productIds).Find(&cats).Error; err != nil {

			return errors.New("unable to find current category")
		}

		// find duplicate
		if len(cats) > 0 {

			return nil
		}

		if err = c.Shadow.Create(&categories).Error; err != nil {

			return errors.New("unable to add category")
		}

		return nil
	}

	// create - linked
	if check, err = c.CreateFast(name, ""); check != nil {

		categories.CategoryId = check.ID

		// create
		if err = c.Shadow.Create(&categories).Error; err != nil {

			return errors.New("unable to add category")
		}

		return nil
	}

	return err
}

func (c *CategoryRepository) Categories(productId uuid.UUID) []string {

	var err error
	var cats []string

	cats = make([]string, 0)

	if bacx.EmptyIdx(productId) {

		return cats
	}

	type Row struct {
		Name string
	}

	ids := bacx.Idx(productId)

	var categories []Row
	categories = make([]Row, 0)

	c.Shadow = c.Shadow.Session(&gorm.Session{})

	prepared := c.Shadow.
		Select("category.name AS name").
		Where("product_id", ids).
		Joins("INNER JOIN category ON categories.category_id = category.id").
		Order("category.name ASC")

	if err = prepared.Error; err != nil {

		return cats
	}

	if err = prepared.Scan(&categories).Error; err != nil {

		return cats
	}

	if len(categories) > 0 {

		for _, category := range categories {

			cats = append(cats, category.Name)
		}
	}

	return cats
}

func (c *CategoryRepository) CatchAll(products []models.Products) []m.KMapImpl {

	c.NewSession()

	data := make([]m.KMapImpl, 0)

	for _, product := range products {

		cats := c.Categories(bacx.Ids(product.ID))

		temp := &m.KMap{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description,
			"image":       pp.Qany(product.Image, nil),
			"price":       product.Price.BigInt(),
			"stocks":      product.Stocks,
			"categories":  cats,
		}

		data = append(data, temp)
	}

	return data
}

func (c *CategoryRepository) UpdateByName(name string, newName string) error {

	if name == "" || newName == "" {

		return errors.New("name is empty string")
	}

	var err error
	var check *models.Category
	if check, _ = c.SearchFast(name); check != nil {

		if err = c.DB.Where("name = ?", name).Updates(map[string]any{

			"name": newName,
		}).Error; err != nil {

			return errors.New("unable to update category")
		}

		return nil
	}

	return errors.New("category not found")
}

func (c *CategoryRepository) DeleteByName(name string) error {

	if name == "" {

		return errors.New("name is empty string")
	}

	var err error
	var check *models.Category
	if check, _ = c.SearchFast(name); check != nil {

		if err = c.Shadow.Unscoped().Where("category_id = ?", check.ID).Delete(&models.Categories{}).Error; err != nil {

			return errors.New("unable to unlink category")
		}

		// unnecessary to delete the category from categories
		//if err = c.DB.Where("name = ?", name).Delete(&models.Category{}).Error; err != nil {
		//
		//	return errors.New("unable to delete category")
		//}

		return nil
	}

	return errors.New("category not found")
}

func (c *CategoryRepository) UnlinkByProductId(productId uuid.UUID) error {

	if bacx.EmptyIdx(productId) {

		return errors.New("invalid productId")
	}

	productIds := bacx.Idx(productId)

	var err error

	if err = c.Shadow.Unscoped().Where("product_id = ?", productIds).Delete(&models.Categories{}).Error; err != nil {

		return errors.New("unable to unlink category")
	}

	return nil
}

func (c *CategoryRepository) NewSession() {

	c.DB = c.DB.Session(&gorm.Session{})
	c.Shadow = c.Shadow.Session(&gorm.Session{})
}
