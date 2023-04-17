package main

import (
	"PapayaNet/app"
	"PapayaNet/papaya"
	"fmt"
)

func main() {

	fmt.Println("Papaya Net v1.0 testing ...")

	pn := papaya.PapayaNet{}

	if err := pn.Init(); err != nil {

		panic(err)
	}

	if err := app.App(&pn); err != nil {

		if err := pn.Close(); err != nil {

			pn.Console.Log(err)
		}
	}

	if err := pn.Close(); err != nil {

		pn.Console.Log(err)
	}
}
