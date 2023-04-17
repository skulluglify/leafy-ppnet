package gen

import (
  "PapayaNet/papaya/koala"
  "reflect"
)

func KMapHunt(value any) bool {

  valueOf := koala.KIndirectValueOf(reflect.ValueOf(value))

  if valueOf.IsValid() {

    switch valueOf.Kind() {
    case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
      return true
    }
  }

  return false
}
