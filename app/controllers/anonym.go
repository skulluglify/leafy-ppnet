package controllers

import (
	"leafy/app/repository"
	"leafy/app/util"
	"net/http"
	"skfw/papaya"
	"skfw/papaya/bunny/swag"
	"skfw/papaya/koala/kornet"
	m "skfw/papaya/koala/mapping"
)

func AnonymController(pn papaya.NetImpl, router swag.SwagRouterImpl) error {

	conn := pn.Connection()
	gorm := conn.GORM()

	productRepo, _ := repository.ProductRepositoryNew(gorm)

	router.Get("/products", &m.KMap{
		"description": "Catch All Products",
		"request": &m.KMap{
			"params": &m.KMap{
				"page": "number",
				"size": "number",
			},
		},
		"responses": swag.OkJSON([]m.KMapImpl{}),
	}, func(ctx *swag.SwagContext) error {

		//////////////////////////////////////

		kReq, _ := ctx.Kornet()

		page := util.ValueToInt(kReq.Query.Get("page"))
		size := util.ValueToInt(kReq.Query.Get("size"))

		var offset int

		if page > 0 {

			offset = page*size - size

			//////////////////////////////////////

			products, err := productRepo.CatchAll(offset, size)

			if err != nil {

				return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew(err.Error(), true))
			}

			var data []m.KMapImpl
			data = make([]m.KMapImpl, 0)

			for _, product := range products {

				data = append(data, &m.KMap{
					"id":          product.ID,
					"name":        product.Name,
					"description": product.Description,
					"price":       product.Price,
					"stocks":      product.Stocks,
				})
			}

			return ctx.Status(http.StatusOK).JSON(data)
		}

		return ctx.Status(http.StatusBadRequest).JSON(kornet.MessageNew("page is zero", true))
	})

	return nil
}
