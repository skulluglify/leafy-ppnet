package main

import (
	"leafy/app"
	"skfw/papaya"
)

func main() {

	pn := papaya.NetNew()

	if err := app.App(pn); err != nil {

		pn.Logger().Error(err)
	}

	if err := pn.Close(); err != nil {

		pn.Logger().Error(err)
	}

	pn.Logger().Log("Shutdown ...")
}
