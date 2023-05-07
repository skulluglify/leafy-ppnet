package test

import (
	"leafy/app/util"
	"testing"
)

func TestAPI(t *testing.T) {

	t.Log(util.NutrientAPI([]string{"orange", "bangkok", "raw"}))
}
