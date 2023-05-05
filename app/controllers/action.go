package controllers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"leafy/app/models"
	"leafy/app/repository"
	"leafy/app/util"
	"net/http"
	"skfw/papaya"
	"skfw/papaya/bunny/swag"
	"skfw/papaya/koala/kornet"
	m "skfw/papaya/koala/mapping"
	mo "skfw/papaya/pigeon/templates/basicAuth/models"
	repo "skfw/papaya/pigeon/templates/basicAuth/repository"
	"time"
)

func ActionController(pn papaya.NetImpl, router swag.SwagRouterImpl) error {

	conn := pn.Connection()
	gorm := conn.GORM()

	userRepo, _ := repository.UserRepositoryNew(gorm)
	productRepo, _ := repository.ProductRepositoryNew(gorm)
	cartRepo, _ := repository.CartRepositoryNew(gorm)
	transactionRepo, _ := repository.TransactionRepositoryNew(gorm)

	router.Get("/info", &m.KMap{
		"AuthToken":   true,
		"description": "User Info",
		"request":     &m.KMap{},
		"responses": swag.OkJSON(&m.KMap{
			"id":           "string",
			"name":         "string",
			"username":     "string",
			"email":        "string",
			"gender":       "string",
			"dob":          "string",
			"address":      "string",
			"phone":        "string",
			"country_code": "string",
			"city":         "string",
			"postal_code":  "string",
			"admin":        "boolean",
			"verify":       "boolean",
			"balance":      "number",
		}),
	}, func(ctx *swag.SwagContext) error {

		var err error

		if ctx.Event() {

			if user, ok := ctx.Target().(*mo.UserModel); ok {

				var balance decimal.Decimal

				idx := repo.Ids(user.ID)

				if balance, err = userRepo.Balance(idx); err != nil {

					return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
				}

				return ctx.Status(http.StatusOK).JSON(&m.KMap{
					"id":           user.ID,
					"name":         user.Name,
					"username":     user.Username,
					"email":        user.Email,
					"gender":       user.Gender,
					"dob":          user.DOB,
					"address":      user.Address,
					"phone":        user.Phone,
					"country_code": user.CountryCode,
					"city":         user.City,
					"postal_code":  user.PostalCode,
					"admin":        user.Admin,
					"verify":       user.Verify,
					"balance":      balance,
				})
			}
		}

		return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to get user information", true))
	})

	router.Put("/info", &m.KMap{
		"AuthToken":   true,
		"description": "User Info",
		"request": &m.KMap{
			"body": swag.JSON(&m.KMap{
				"name":         "string",
				"gender":       "string",
				"dob":          "string",
				"address":      "string",
				"phone":        "string",
				"country_code": "string",
				"city":         "string",
				"postal_code":  "string",
			}),
		},
		"responses": swag.OkJSON(&kornet.Message{}),
	}, func(ctx *swag.SwagContext) error {

		if ctx.Event() {

			if user, ok := ctx.Target().(*mo.UserModel); ok {

				kReq, _ := ctx.Kornet()

				data := &m.KMap{}

				if err := json.Unmarshal(kReq.Body.ReadAll(), data); err != nil {

					return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to parse json", true))
				}

				name := m.KValueToString(data.Get("name"))
				gender := m.KValueToString(data.Get("gender"))
				dob := m.KValueToString(data.Get("dob"))
				address := m.KValueToString(data.Get("address"))
				phone := m.KValueToString(data.Get("phone"))
				countryCode := m.KValueToString(data.Get("country_code"))
				city := m.KValueToString(data.Get("city"))
				postalCode := m.KValueToString(data.Get("postal_code"))

				// 0001-01-01 00:00:00

				var err error
				var DOB time.Time

				if DOB, err = time.Parse("2006-01-02 15:04:05", dob); err != nil {

					return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew("unable to parse date of birth", true))
				}

				userNew := &models.Users{
					UserModel: &mo.UserModel{
						ID:          user.ID,
						Name:        name,
						Username:    user.Username,
						Email:       user.Email,
						Gender:      gender,
						DOB:         DOB,
						Address:     address,
						Phone:       phone,
						CountryCode: countryCode,
						City:        city,
						PostalCode:  postalCode,
					},
				}

				if err = userRepo.Update(userNew); err != nil {

					return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
				}

				return ctx.Status(http.StatusOK).JSON(kornet.MessageNew("update user information", false))
			}
		}

		return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to get user information", true))
	})

	router.Get("/carts", &m.KMap{
		"AuthToken":   true,
		"description": "Catch All Carts",
		"request": &m.KMap{
			"params": &m.KMap{
				"page": "number",
				"size": "number",
			},
		},
		"responses": swag.OkJSON([]m.KMapImpl{}),
	}, func(ctx *swag.SwagContext) error {

		var err error
		var carts []models.Cart

		if ctx.Event() {

			if user, ok := ctx.Target().(*mo.UserModel); ok {

				kReq, _ := ctx.Kornet()

				page := util.ValueToInt(kReq.Query.Get("page"))
				size := util.ValueToInt(kReq.Query.Get("size"))

				var offset int

				if page > 0 {

					offset = page*size - size

					idx := repo.Ids(user.ID)

					if carts, err = cartRepo.CatchAll(idx, offset, size); err != nil {

						return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
					}

					var data []m.KMapImpl
					data = make([]m.KMapImpl, 0)

					for _, cart := range carts {

						var product *models.Products

						if product, err = productRepo.SearchUnscopedFastById(repo.Ids(cart.ProductID)); err != nil {

							return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
						}

						data = append(data, &m.KMap{
							"id": cart.ID,
							"product": &m.KMap{
								"id":     product.ID,
								"name":   product.Name,
								"stocks": product.Stocks,
							},
							"qty": cart.Qty,
						})
					}

					return ctx.Status(http.StatusOK).JSON(data)
				}

				return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew("page is zero", true))
			}
		}

		return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to get user information", true))
	})

	router.Post("/cart", &m.KMap{
		"AuthToken":   true,
		"description": "Create New Cart",
		"request": &m.KMap{
			"body": swag.JSON(&m.KMap{
				"productId": "string",
				"qty":       "number",
			}),
		},
		"responses": swag.CreatedJSON(&kornet.Message{}),
	}, func(ctx *swag.SwagContext) error {

		if ctx.Event() {

			if user, ok := ctx.Target().(*mo.UserModel); ok {

				kReq, _ := ctx.Kornet()

				data := &m.KMap{}

				if err := json.Unmarshal(kReq.Body.ReadAll(), data); err != nil {

					return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to parse json", true))
				}

				productId := m.KValueToString(data.Get("productId"))
				qty := util.ValueToInt(data.Get("qty"))

				userIdx := repo.Ids(user.ID)
				productIdx := repo.Ids(productId)

				// check stocks
				if check, _ := productRepo.SearchFastById(productIdx); check != nil {

					if check.Stocks < qty {

						return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew("quantity out of stocks", true))
					}
				}

				if _, err := cartRepo.CreateFast(userIdx, productIdx, qty); err != nil {

					return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
				}

				return ctx.Status(http.StatusCreated).JSON(kornet.MessageNew("create new cart", false))
			}
		}

		return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to get user information", true))
	})

	router.Put("/cart", &m.KMap{
		"AuthToken":   true,
		"description": "Update Cart",
		"request": &m.KMap{
			"params": &m.KMap{
				"id": "string",
			},
			"body": swag.JSON(&m.KMap{
				"qty": "number",
			}),
		},
	}, func(ctx *swag.SwagContext) error {

		if ctx.Event() {

			if user, ok := ctx.Target().(*mo.UserModel); ok {

				kReq, _ := ctx.Kornet()

				ids := m.KValueToString(kReq.Query.Get("id"))

				idx := repo.Ids(ids)

				data := &m.KMap{}

				if err := json.Unmarshal(kReq.Body.ReadAll(), data); err != nil {

					return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to parse json", true))
				}

				qty := util.ValueToInt(data.Get("qty"))

				userIdx := repo.Ids(user.ID)

				cart := &models.Cart{
					ID:  ids,
					Qty: qty,
				}

				var productIdx uuid.UUID

				oldCart, err := cartRepo.SafeSearchFastById(idx, userIdx)

				if err != nil {

					return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to get cart information", true))
				}

				productId := oldCart.ProductID
				productIdx = repo.Ids(productId)

				// check stocks
				if check, _ := productRepo.SearchFastById(productIdx); check != nil {

					if check.Stocks < qty {

						return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew("quantity out of stocks", true))
					}
				}

				if err := cartRepo.SafeUpdate(userIdx, cart); err != nil {

					return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to update cart", true))
				}

				return ctx.Status(http.StatusOK).JSON(kornet.MessageNew("update cart", false))
			}
		}

		return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to get user information", true))
	})

	router.Delete("/cart", &m.KMap{
		"AuthToken":   true,
		"description": "Delete Cart",
		"request": &m.KMap{
			"params": &m.KMap{
				"id": "string",
			},
		},
		"responses": swag.OkJSON(&kornet.Message{}),
	}, func(ctx *swag.SwagContext) error {

		if ctx.Event() {

			if user, ok := ctx.Target().(*mo.UserModel); ok {

				kReq, _ := ctx.Kornet()

				ids := m.KValueToString(kReq.Query.Get("id"))

				idx := repo.Ids(ids)
				userIdx := repo.Ids(user.ID)

				if err := cartRepo.SafeDelete(idx, userIdx); err != nil {

					return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to delete cart", true))
				}

				return ctx.Status(http.StatusOK).JSON(kornet.MessageNew("delete cart", false))
			}
		}

		return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to get user information", true))
	})

	router.Get("/transactions", &m.KMap{
		"AuthToken":   true,
		"description": "Catch All Transactions",
		"request": &m.KMap{
			"params": &m.KMap{
				"page": "number",
				"size": "number",
			},
		},
		"responses": swag.OkJSON([]m.KMapImpl{}),
	}, func(ctx *swag.SwagContext) error {

		var err error
		var transactions []models.Transactions

		if ctx.Event() {

			if user, ok := ctx.Target().(*mo.UserModel); ok {

				kReq, _ := ctx.Kornet()

				page := util.ValueToInt(kReq.Query.Get("page"))
				size := util.ValueToInt(kReq.Query.Get("size"))

				var offset int

				if page > 0 {

					offset = page*size - size

					idx := repo.Ids(user.ID)

					if transactions, err = transactionRepo.CatchAll(idx, offset, size); err != nil {

						return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
					}

					var data []m.KMapImpl
					data = make([]m.KMapImpl, 0)

					// normalize data object for transactions map safety

					for _, transaction := range transactions {

						var carts []m.KMapImpl
						carts = make([]m.KMapImpl, 0)

						for _, cart := range transaction.Carts {

							var product *models.Products

							productId := repo.Ids(cart.ProductID)

							if product, err = productRepo.SearchFastById(productId); err != nil {

								return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
							}

							carts = append(carts, &m.KMap{
								"id": cart.ID,
								"product": &m.KMap{
									"id":   product.ID,
									"name": product.Name,
								},
								"qty": cart.Qty,
							})
						}

						data = append(data, &m.KMap{
							"id":             transaction.ID,
							"carts":          carts,
							"payment_method": transaction.PaymentMethod,
							"verify":         transaction.Verify,
						})
					}

					return ctx.Status(http.StatusOK).JSON(data)
				}

				return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew("page is zero", true))
			}
		}

		return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to get user information", true))
	})

	router.Post("/transaction", &m.KMap{
		"AuthToken":   true,
		"description": "Create New Transaction",
		"request": &m.KMap{
			"body": swag.JSON(&m.KMap{
				"payment_method": "string",
			}),
		},
		"responses": swag.OkJSON(&kornet.Message{}),
	}, func(ctx *swag.SwagContext) error {

		var err error

		if ctx.Event() {

			if user, ok := ctx.Target().(*mo.UserModel); ok {

				kReq, _ := ctx.Kornet()

				data := &m.KMap{}

				if err = json.Unmarshal(kReq.Body.ReadAll(), data); err != nil {

					return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to parse json", true))
				}

				paymentMethod := m.KValueToString(data.Get("payment_method"))

				userIdx := repo.Ids(user.ID)

				if err = userRepo.PayCheck(userIdx); err != nil {

					return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew(err.Error(), true))
				}

				var transaction *models.Transactions

				if transaction, err = transactionRepo.CreateFast(userIdx, paymentMethod); err != nil {

					return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
				}

				if err = transactionRepo.Verify(transaction); err != nil {

					return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
				}

				return ctx.Status(http.StatusCreated).JSON(kornet.MessageNew("create new transaction", false))
			}
		}

		return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to get user information", true))
	})

	router.Get("/bill", &m.KMap{
		"AuthToken":   true,
		"description": "Show BIll",
		"responses":   swag.OkJSON([]m.KMapImpl{}),
	}, func(ctx *swag.SwagContext) error {

		var err error
		var bill decimal.Decimal

		if ctx.Event() {

			if user, ok := ctx.Target().(*mo.UserModel); ok {

				idx := repo.Ids(user.ID)

				if bill, err = userRepo.Bill(idx); err != nil {

					return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
				}

				return ctx.Status(http.StatusOK).JSON(&m.KMap{
					"pay": bill.BigInt(),
				})
			}
		}

		return nil
	})

	return nil
}
