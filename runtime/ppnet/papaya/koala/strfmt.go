package koala

import (
	"PapayaNet/papaya/panda"
	"strconv"
)

// Method String to Number

func KStrToNum(value string) int {

	// try convert
	if v, err := strconv.Atoi(value); err == nil {

		return v
	}

	// default value
	return 0
}

// Method convert String to Boolean

func KStrToBool(value string) bool {

	// try convert
	if v, err := strconv.ParseBool(value); err == nil {

		return v
	}

	// default value
	return false
}

func KStrZeroFill(text string, s int) string {

	var zeros string

	k := panda.Min(len(text), s)
	z := s - k

	for i := 0; i < z; i++ {

		zeros += "0"
	}

	return zeros + text[:k]
}

func KStrPadStart(text string, s int) string {

	var pads string

	k := panda.Min(len(text), s)
	z := s - k

	for i := 0; i < z; i++ {

		pads += " "
	}

	return pads + text[:k]
}

func KStrPadEnd(text string, s int) string {

	var pads string

	k := panda.Min(len(text), s)
	z := s - k

	for i := 0; i < z; i++ {

		pads += " "
	}

	return text[:k] + pads
}
