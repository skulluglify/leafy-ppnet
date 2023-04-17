package main

import (
	"PapayaNet/papaya/koala"
	"PapayaNet/papaya/koala/collection"
)

func main() {

	var i uint
	console := koala.KConsoleNew()
	console.Log("KList testing ...")

	list := collection.KListNewR[int]([]int{12, 24, 36, 48, 60, 72, 84, 96, 108})

	for i = 0; i < list.Len(); i++ {

		console.Log(list.Get(i))
	}

	console.Warn("splice ...")

	// TODO: fix problem --
	// TODO: initial val.Type() after val.IsValid() ++
	removes := list.Splice(1, 1, 72, 80)

	for i = 0; i < removes.Len(); i++ {

		console.Log(removes.Get(i))
	}

	console.Warn("look ...")
	console.Warn(list.Len())

	for i = 0; i < list.Len(); i++ {

		console.Log(list.Get(i))
	}
}
