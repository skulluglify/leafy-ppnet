package app

import (
	"skfw/papaya"
	"skfw/papaya/pigeon/cors"
)

func ManageControlResourceShared(pn papaya.NetImpl) error {

	manageConsumers, err := cors.ManageConsumersNew()
	if err != nil {

		return err
	}

	//manageConsumers.Grant("*")
	manageConsumers.Grant("http://localhost")
	manageConsumers.Grant("http://localhost:8000")
	manageConsumers.Grant("http://159.89.202.201")
	manageConsumers.Grant("http://skfw.net")

	pn.Use(cors.MakeMiddlewareForManageConsumers(manageConsumers))

	return nil
}
