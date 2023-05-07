package repository

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"leafy/app/models"
	"leafy/app/util"
	m "skfw/papaya/koala/mapping"
	"skfw/papaya/pigeon/templates/basicAuth/repository"
)

type NutrientRepository struct {
	DB *gorm.DB
}

type NutrientRepositoryImpl interface {
	Init(DB *gorm.DB) error
	SearchFastByProductId(productId uuid.UUID) ([]models.Nutrients, error)
	MapCatchAll(data []m.KMapImpl) []m.KMapImpl
	DeleteFast(productId uuid.UUID) error
	CreateFast(productId uuid.UUID, categories []string) ([]models.Nutrients, error)
}

func NutrientRepositoryNew(DB *gorm.DB) (NutrientRepositoryImpl, error) {

	nutrientRepo := &NutrientRepository{}
	if err := nutrientRepo.Init(DB); err != nil {

		return nil, err
	}
	return nutrientRepo, nil
}

func (n *NutrientRepository) Init(DB *gorm.DB) error {

	if DB != nil {

		n.DB = DB.Model(&models.Nutrients{})

		return nil
	}

	return errors.New("DB is NULL")
}

func (n *NutrientRepository) SearchFastByProductId(productId uuid.UUID) ([]models.Nutrients, error) {

	var err error

	n.NewSession()

	if repository.EmptyIdx(productId) {

		return nil, errors.New("invalid productId")
	}

	productIds := repository.Idx(productId)

	var nutrients []models.Nutrients
	nutrients = make([]models.Nutrients, 0)

	if err = n.DB.Where("product_id = ?", productIds).Find(&nutrients).Error; err != nil {

		return nil, errors.New("unable to find nutrient")
	}

	return nutrients, nil
}

func (n *NutrientRepository) MapCatchAll(data []m.KMapImpl) []m.KMapImpl {

	var err error
	var nutrients []models.Nutrients

	for _, product := range data {

		ids := m.KValueToString(product.Get("id"))
		id := repository.Ids(ids)

		if repository.EmptyIdx(id) {

			continue
		}

		if nutrients, err = n.SearchFastByProductId(id); err != nil {

			continue
		}

		var temp []m.KMapImpl
		temp = make([]m.KMapImpl, 0)

		if len(nutrients) > 0 {

			for _, nutrient := range nutrients {

				temp = append(temp, &m.KMap{
					"name":  nutrient.Name,
					"unit":  nutrient.Unit,
					"value": nutrient.Value,
				})
			}

			product.Put("nutrients", temp)
		}
	}

	return data
}

func (n *NutrientRepository) DeleteFast(productId uuid.UUID) error {

	var err error

	if repository.EmptyIdx(productId) {

		return errors.New("invalid productId")
	}

	productIds := repository.Idx(productId)

	if err = n.DB.Unscoped().Where("product_id = ?", productIds).Delete(&models.Nutrients{}).Error; err != nil {

		return errors.New("unable to delete datasheet nutrients")
	}

	return nil
}

func (n *NutrientRepository) CreateFast(productId uuid.UUID, categories []string) ([]models.Nutrients, error) {

	var err error

	var nutrients []models.Nutrients
	nutrients = make([]models.Nutrients, 0)

	if nutrients, err = n.SearchFastByProductId(productId); len(nutrients) > 0 {

		return nutrients, errors.New("nutrients has been added")
	}

	productIds := repository.Idx(productId)

	// find nutrients
	var data []util.Nutrient
	data = util.NutrientAPI(categories)

	if len(data) > 0 {

		for _, nutrient := range data {

			nut := models.Nutrients{
				ProductId: productIds,
				Name:      nutrient.Name,
				Unit:      nutrient.Unit,
				Value:     nutrient.Value,
			}

			// pass - error
			if err = n.DB.Create(&nut).Error; err != nil {

				continue
			}

			nutrients = append(nutrients, nut)
		}

		return nutrients, nil
	}

	return nil, errors.New("unable to search nutrients")
}

func (n *NutrientRepository) NewSession() {

	n.DB = n.DB.Session(&gorm.Session{})
}
