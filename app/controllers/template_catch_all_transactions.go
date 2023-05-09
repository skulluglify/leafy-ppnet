package controllers

import (
	"errors"
	"github.com/google/uuid"
	"leafy/app/models"
	"leafy/app/repository"
	"skfw/papaya"
	m "skfw/papaya/koala/mapping"
	bacx "skfw/papaya/pigeon/templates/basicAuth/util"
)

type CatchAllTransactionsHandler func(userId uuid.UUID, page int, size int) ([]m.KMapImpl, error)

func TemplateCatchAllTransactions(pn papaya.NetImpl) CatchAllTransactionsHandler {

	conn := pn.Connection()
	gorm := conn.GORM()

	productRepo, _ := repository.ProductRepositoryNew(gorm)
	transactionRepo, _ := repository.TransactionRepositoryNew(gorm)

	return func(userId uuid.UUID, page int, size int) ([]m.KMapImpl, error) {

		var err error

		var offset int

		var data []m.KMapImpl

		if page > 0 {

			data = make([]m.KMapImpl, 0)

			offset = page*size - size

			transactions := make([]models.Transactions, 0)

			if transactions, err = transactionRepo.CatchAll(userId, offset, size); err != nil {

				return data, err
			}

			// normalize data object for transactions map safety

			for _, transaction := range transactions {

				var carts []m.KMapImpl
				carts = make([]m.KMapImpl, 0)

				for _, cart := range transaction.Carts {

					var product *models.Products

					productId := bacx.Ids(cart.ProductID)

					if product, err = productRepo.SearchFastById(productId); err != nil {

						return data, err
					}

					carts = append(carts, &m.KMap{
						"id": cart.ID,
						"product": &m.KMap{
							"id":     product.ID,
							"name":   product.Name,
							"stocks": product.Stocks,
						},
						"qty": cart.Qty,
					})
				}

				data = append(data, &m.KMap{
					"id":             transaction.ID,
					"carts":          carts,
					"payment_method": transaction.PaymentMethod,
				})
			}

			return data, nil
		}

		return nil, errors.New("page is zero")
	}
}
