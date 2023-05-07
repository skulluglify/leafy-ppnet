package controllers

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"leafy/app/models"
	"leafy/app/repository"
	"leafy/app/util"
	"mime"
	"mime/multipart"
	"net/http"
	"skfw/papaya"
	"skfw/papaya/bunny/swag"
	"skfw/papaya/koala/kornet"
	m "skfw/papaya/koala/mapping"
	"skfw/papaya/koala/pp"
	repo "skfw/papaya/pigeon/templates/basicAuth/repository"
)

func AdminController(pn papaya.NetImpl, router swag.SwagRouterImpl) error {

	conn := pn.Connection()
	gorm := conn.GORM()

	userRepo, _ := repository.UserRepositoryNew(gorm)
	productRepo, _ := repository.ProductRepositoryNew(gorm)
	categoryRepo, _ := repository.CategoryRepositoryNew(gorm)
	nutrientRepo, _ := repository.NutrientRepositoryNew(gorm)

	catchAllTransactions := TemplateCatchAllTransactions(pn)

	router.Post("/topup", &m.KMap{
		"AuthToken":   true,
		"Admin":       true,
		"description": "Topup Money",
		"request": &m.KMap{
			"params": &m.KMap{
				"userId": "string",
			},
			"body": swag.JSON(&m.KMap{
				"balance": "number",
			}),
		},
		"responses": swag.OkJSON(&kornet.Message{}),
	}, func(ctx *swag.SwagContext) error {

		var err error

		kReq, _ := ctx.Kornet()

		data := m.KMap{}

		if err = json.Unmarshal(kReq.Body.ReadAll(), &data); err != nil {

			return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to parse json", true))
		}

		balance := decimal.NewFromInt(int64(util.ValueToInt(data.Get("balance"))))

		userId := m.KValueToString(kReq.Query.Get("userId"))

		userIdx := repo.Ids(userId)

		if err = userRepo.Topup(userIdx, balance); err != nil {

			return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
		}

		return ctx.Status(http.StatusOK).JSON(kornet.MessageNew("successful topup user balance", false))
	})

	router.Get("/sessions", &m.KMap{
		"AuthToken":   true,
		"Admin":       true,
		"description": "Catch All Sessions",
		"request": &m.KMap{
			"params": &m.KMap{
				"page": "number",
				"size": "number",
			},
		},
		"responses": swag.OkJSON([]m.KMapImpl{}),
	}, func(ctx *swag.SwagContext) error {

		kReq, _ := ctx.Kornet()

		page := util.ValueToInt(kReq.Query.Get("page"))
		size := util.ValueToInt(kReq.Query.Get("size"))

		var offset int

		if page > 0 {

			offset = page*size - size

			users, err := userRepo.CatchAll(offset, size)
			if err != nil {

				return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
			}

			var data []m.KMapImpl
			var sessions []m.KMapImpl

			data = make([]m.KMapImpl, 0)

			for _, user := range users {

				sessions = make([]m.KMapImpl, 0)

				for _, s := range user.Sessions {

					// look up non deleted
					if !s.DeletedAt.Valid {

						sessions = append(sessions, &m.KMap{
							"id":             s.ID,
							"client_ip":      s.ClientIP,
							"user_agent":     s.UserAgent,
							"expired":        s.Expired,
							"last_activated": s.LastActivated,
						})
					}
				}

				data = append(data, &m.KMap{
					"username": user.Username,
					"name":     user.Name,
					"email":    user.Email,
					"phone":    user.Phone,
					"sessions": sessions,
				})
			}

			return ctx.Status(http.StatusOK).JSON(data)
		}

		return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew("page is zero", true))
	})

	router.Post("/product", &m.KMap{
		"AuthToken":   true,
		"Admin":       true,
		"description": "Create New Product",
		"request": &m.KMap{
			"body": swag.JSON(&m.KMap{
				"name":        "string",
				"description": "string",
				"price":       "number",
				"stocks":      "number",
				"categories":  []string{},
			}),
		},
		"responses": swag.CreatedJSON(&kornet.Message{}),
	}, func(ctx *swag.SwagContext) error {

		var err error

		kReq, _ := ctx.Kornet()

		data := &m.KMap{}

		if err = json.Unmarshal(kReq.Body.ReadAll(), data); err != nil {

			return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to parse data", true))
		}

		name := m.KValueToString(data.Get("name"))
		description := m.KValueToString(data.Get("description"))
		price := decimal.NewFromFloat(m.KValueToFloat(data.Get("price")))
		stocks := util.ValueToInt(data.Get("stocks"))
		categories := util.ValueToArrayStr(data.Get("categories"))

		var check *models.Products

		if check, err = productRepo.CreateFast(name, description, price, stocks); err != nil {

			return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
		}

		// add category
		if check != nil {

			productId := repo.Ids(check.ID)

			for _, cate := range categories {

				if err = categoryRepo.Add(productId, cate); err != nil {

					return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
				}
			}

			if _, err = nutrientRepo.CreateFast(productId, categories); err != nil {

				return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
			}
		}

		return ctx.Status(http.StatusCreated).JSON(kornet.MessageNew("create new product", false))
	})

	router.Post("/category", &m.KMap{
		"AuthToken":   true,
		"Admin":       true,
		"description": "Add Category",
		"request": &m.KMap{
			"params": &m.KMap{
				"productId": "string",
			},
			"body": swag.JSON(&m.KMap{
				"name": "string",
			}),
		},
		"responses": swag.OkJSON(&kornet.Message{}),
	}, func(ctx *swag.SwagContext) error {

		var err error

		kReq, _ := ctx.Kornet()

		data := &m.KMap{}

		if err = json.Unmarshal(kReq.Body.ReadAll(), data); err != nil {

			return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to parse data", true))
		}

		productId := m.KValueToString(kReq.Query.Get("productId"))
		name := m.KValueToString(data.Get("name"))
		//description := m.KValueToString(data.Get("description"))

		productIdx := repo.Ids(productId)

		if check, _ := productRepo.SearchFastById(productIdx); check != nil {

			// try to create a new category
			//if _, err = categoryRepo.SearchFast(name); err != nil {
			//
			//	if _, err = categoryRepo.CreateFast(name, description); err != nil {
			//
			//		return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
			//	}
			//}

			if err = categoryRepo.Add(productIdx, name); err != nil {

				return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
			}

			return ctx.Status(http.StatusOK).JSON(kornet.MessageNew("link category", false))
		}

		return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew("product not found", false))
	})

	router.Put("/category", &m.KMap{
		"AuthToken":   true,
		"Admin":       true,
		"description": "Update Category",
		"request": &m.KMap{
			"params": &m.KMap{
				"name": "string",
			},
			"body": swag.JSON(&m.KMap{
				"name": "string",
			}),
		},
		"responses": swag.OkJSON(&kornet.Message{}),
	}, func(ctx *swag.SwagContext) error {

		var err error

		kReq, _ := ctx.Kornet()

		data := &m.KMap{}

		if err = json.Unmarshal(kReq.Body.ReadAll(), data); err != nil {

			return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to parse data", true))
		}

		name := m.KValueToString(kReq.Query.Get("name"))
		newName := m.KValueToString(data.Get("name"))

		if err = categoryRepo.UpdateByName(name, newName); err != nil {

			return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
		}

		return ctx.Status(http.StatusOK).JSON(kornet.MessageNew("update category", false))
	})

	router.Delete("/category", &m.KMap{
		"AuthToken":   true,
		"Admin":       true,
		"description": "Delete Category",
		"request": &m.KMap{
			"params": &m.KMap{
				"name": "string",
			},
		},
		"responses": swag.OkJSON(&kornet.Message{}),
	}, func(ctx *swag.SwagContext) error {

		var err error

		kReq, _ := ctx.Kornet()

		name := m.KValueToString(kReq.Query.Get("name"))

		if err = categoryRepo.DeleteByName(name); err != nil {

			return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
		}

		return ctx.Status(http.StatusOK).JSON(kornet.MessageNew("delete category", false))
	})

	router.Put("/product", &m.KMap{
		"AuthToken":   true,
		"Admin":       true,
		"description": "Update Product",
		"request": &m.KMap{
			"params": &m.KMap{
				"id": "string",
			},
			"body": swag.JSON(&m.KMap{
				"name":        "string",
				"description": "string",
				"price":       "number",
				"stocks":      "number",
				"categories":  []string{},
			}),
		},
		"responses": swag.CreatedJSON(&kornet.Message{}),
	}, func(ctx *swag.SwagContext) error {

		var err error

		kReq, _ := ctx.Kornet()

		ids := m.KValueToString(kReq.Query.Get("id"))

		data := &m.KMap{}

		if err := json.Unmarshal(kReq.Body.ReadAll(), data); err != nil {

			return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to parse data", true))
		}

		name := m.KValueToString(data.Get("name"))
		description := m.KValueToString(data.Get("description"))
		price := decimal.NewFromFloat(m.KValueToFloat(data.Get("price")))
		stocks := util.ValueToInt(data.Get("stocks"))
		categories := util.ValueToArrayStr(data.Get("categories"))

		id := repo.Ids(ids)

		if !repo.EmptyIdx(id) {

			if product, _ := productRepo.SearchFastById(id); product != nil {

				product.Name = name
				product.Description = description
				product.Price = price
				product.Stocks = stocks

				if productRepo.Update(id, product) != nil {

					return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to update product", true))
				}

				for _, cate := range categories {

					if err = categoryRepo.Add(repo.Ids(product.ID), cate); err != nil {

						return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
					}
				}

				return ctx.Status(http.StatusOK).JSON(kornet.MessageNew("update product", false))
			}

			return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew("product not found", true))
		}

		return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew("unable to parse id", true))

	})

	router.Delete("/product", &m.KMap{
		"AuthToken":   true,
		"Admin":       true,
		"description": "Delete Product",
		"request": &m.KMap{
			"params": &m.KMap{
				"id": "string",
			},
		},
		"responses": swag.CreatedJSON(&kornet.Message{}),
	}, func(ctx *swag.SwagContext) error {

		var err error

		kReq, _ := ctx.Kornet()

		query := kReq.Query

		ids := m.KValueToString(query.Get("id"))

		id := repo.Ids(ids)

		if repo.EmptyIdx(id) {

			return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew("id is empty", true))
		}

		if err = nutrientRepo.DeleteFast(id); err != nil {

			return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
		}

		if err = categoryRepo.UnlinkByProductId(id); err != nil {

			return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
		}

		if err = productRepo.DeleteFast(id); err != nil {

			return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
		}

		return ctx.Status(http.StatusOK).JSON(kornet.MessageNew("delete product", false))

	})

	router.Get("/transactions", &m.KMap{
		"AuthToken":   true,
		"Admin":       true,
		"description": "Catch All Transactions",
		"request": &m.KMap{
			"params": &m.KMap{
				"page":      "number",
				"size":      "number",
				"maxCatch?": "number",
			},
		},
		"responses": swag.OkJSON([]m.KMapImpl{}),
	}, func(ctx *swag.SwagContext) error {

		var data []m.KMapImpl

		kReq, _ := ctx.Kornet()

		page := util.ValueToInt(kReq.Query.Get("page"))
		size := util.ValueToInt(kReq.Query.Get("size"))
		maxCatch := pp.QInt(util.ValueToInt(kReq.Query.Get("maxCatch")), 200)

		var offset int

		if page > 0 {

			data = make([]m.KMapImpl, 0)

			offset = page*size - size

			users, err := userRepo.CatchAll(offset, size)

			if err != nil {
				return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
			}

			for _, user := range users {

				userId := repo.Ids(user.ID)

				var transactions []m.KMapImpl

				if transactions, err = catchAllTransactions(userId, 1, maxCatch); err != nil {

					return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
				}

				data = append(data, &m.KMap{
					"user": &m.KMap{
						"id":       user.ID,
						"name":     user.Name,
						"username": user.Username,
						"email":    user.Email,
						"phone":    user.Phone,
					},
					"transactions": transactions,
				})
			}

			return ctx.Status(http.StatusOK).JSON(data)
		}

		return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew("page is zero", true))
	})

	router.Post("/product/upload", &m.KMap{
		"AuthToken":   true,
		"Admin":       true,
		"description": "Upload Product Image",
		"request": &m.KMap{
			"params": &m.KMap{
				"productId": "string",
			},
		},
		"responses": swag.CreatedJSON(&kornet.Message{}),
	}, func(ctx *swag.SwagContext) error {

		var err error

		kReq, _ := ctx.Kornet()

		productIds := m.KValueToString(kReq.Query.Get("productId"))

		if productIds == "" {

			return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew("invalid productId", true))
		}

		productId := repo.Ids(productIds)

		var check *models.Products

		if check, err = productRepo.SearchFastById(productId); err != nil {

			return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew("product not found", true))
		}

		name := util.SafePathName(check.Name)

		var form *multipart.Form
		var extensions []string
		var ext string

		if form, err = ctx.MultipartForm(); err != nil {

			return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew("bad request", true))
		}

		images := m.Keys([]string{
			"image/jpeg",
			"image/png",
		})

		var found bool

		found = false

		for k, h := range form.File {

			switch k {
			case "img", "image":

				if len(h) > 0 {

					header := h[0]
					cTy := header.Header.Get("Content-Type")

					if images.Contain(cTy) {

						if extensions, err = mime.ExtensionsByType(cTy); err != nil {

							return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to get name of extension", true))
						}

						n := len(extensions)

						if n > 0 {

							ext = extensions[n-1] // the last thing maybe a good choice

							ext = ".png" // force use PNG formatter

							output := "assets/public/products/" + name + ext

							//if err = ctx.SaveFile(header, "assets/public/products/"+name+ext); err != nil {
							//
							//	return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to save image", true))
							//}

							check.Image = name + ext

							if err = productRepo.Update(productId, check); err != nil {

								return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
							}

							return util.SaveImage(k, ext, output)(ctx.Ctx)
						}
					}
				}

				found = true
				break
			}

			if found {

				break
			}
		}

		if !found {

			return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew("unsupported format", true))
		}

		return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("something wrong", true))
	})

	return nil
}
