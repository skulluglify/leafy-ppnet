package app

import (
	"leafy/app/controllers"
	"leafy/app/models"
	"skfw/papaya"
	"skfw/papaya/bunny/swag"
	"skfw/papaya/pigeon/templates/basicAuth/repository"
	"time"
)

func App(pn papaya.NetImpl) error {

	conn := pn.Connection()
	gorm := conn.GORM()

	ManageControlResourceShared(pn)

	swagger := pn.MakeSwagger(&swag.SwagInfo{
		Title:       "Leafy API",
		Version:     "1.0.0",
		Description: "Leafy Marketplace API for Documentation",
	})

	mainGroup := swagger.Group("/api/v1", "Schema")
	userGroup := mainGroup.Group("/users", "Authentication")

	userRouter := userGroup.Router()

	expired := time.Hour * 4
	activeDuration := time.Minute * 30 // interval
	maxSessions := 6

	basicAuth := repository.BasicAuthNew(conn, expired, activeDuration, maxSessions)
	basicAuth.Bind(swagger, userRouter)

	gorm.AutoMigrate(&models.User{}, &models.Session{}, &models.Product{}, &models.Cart{}, &models.Transaction{})

	controllers.UserController(pn, userRouter)

	swagger.Start()

	return pn.Serve("127.0.0.1", 8000)
}
