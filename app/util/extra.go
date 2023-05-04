package util

import (
	"reflect"
	"skfw/papaya/koala/mapping"
	"skfw/papaya/koala/pp"
)

func ValueToInt(value any) int {

	val := pp.KIndirectValueOf(value)

	if val.IsValid() {

		ty := val.Type()

		switch ty.Kind() {

		case reflect.Float64:

			return int(mapping.KValueToFloat(value))
		}

		return int(mapping.KValueToInt(value))
	}

	return 0
}
