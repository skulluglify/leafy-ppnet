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

func ValueToArrayStr(data any) []string {

	var temp []string
	temp = make([]string, 0)

	val := pp.KIndirectValueOf(data)

	if val.IsValid() {

		ty := val.Type()

		switch ty.Kind() {
		case reflect.Array, reflect.Slice:

			// loop - validity - casting

			for i := 0; i < val.Len(); i++ {

				elem := val.Index(i)

				vElem := pp.KIndirectValueOf(elem)

				if vElem.IsValid() {

					tyElem := vElem.Type()

					switch tyElem.Kind() {

					case reflect.String:

						temp = append(temp, vElem.String())
					}
				}
			}
		}
	}

	return temp
}
