package app

import (
	"leafy/app/controllers"
	"leafy/app/models"
	"leafy/app/tasks"
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

	anonymGroup := mainGroup.Group("/", "Anonymous")
	userGroup := mainGroup.Group("/users", "Authentication")
	actionGroup := mainGroup.Group("/users", "Action")
	adminGroup := mainGroup.Group("/admin", "Administrator")

	anonymRouter := anonymGroup.Router()
	userRouter := userGroup.Router()
	actionRouter := actionGroup.Router()
	adminRouter := adminGroup.Router()

	controllers.AnonymController(pn, anonymRouter)

	expired := time.Hour * 4
	activeDuration := time.Minute * 30 // interval
	maxSessions := 6

	basicAuth := repository.BasicAuthNew(conn, expired, activeDuration, maxSessions)
	basicAuth.Bind(swagger, userRouter)

	swagger.AddTask(tasks.MakeAdminTask())

	gorm.AutoMigrate(
		&models.Users{},
		&models.Sessions{},
		&models.Products{},
		&models.Cart{},
		&models.Transactions{},
		&models.Category{},
		&models.Categories{},
		&models.Nutrients{},
	)

	controllers.ActionController(pn, actionRouter)
	controllers.AdminController(pn, adminRouter)

	swagger.Start()

	return pn.Serve("127.0.0.1", 8000)
}
