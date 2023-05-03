package controllers

import (
	"net/http"
	"skfw/papaya"
	"skfw/papaya/bunny/swag"
	"skfw/papaya/koala/kornet"
	m "skfw/papaya/koala/mapping"
	mo "skfw/papaya/pigeon/templates/basicAuth/models"
)

func UserController(pn papaya.NetImpl, router swag.SwagRouterImpl) error {

	//conn := pn.Connection()
	//GORM := conn.GORM()
	//
	//userRepo, _ := repository.UserRepositoryNew(GORM)

	router.Get("/info", &m.KMap{
		"AuthToken":   true,
		"description": "User Info",
		"request":     &m.KMap{},
		"responses": swag.OkJSON(&m.KMap{
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
		}),
	}, func(ctx *swag.SwagContext) error {

		if ctx.Event() {

			if user, ok := ctx.Target().(*mo.UserModel); ok {

				return ctx.Status(http.StatusOK).JSON(&m.KMap{
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
				})
			}
		}

		return ctx.Status(http.StatusInternalServerError).JSON(kornet.MessageNew("unable to get user information", true))
	})

	return nil
}
