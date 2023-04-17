package koala

import "reflect"

// hack valueOf from `reflect.Value` to get indirect as value
// read as interface ptr, ptr, and passing
// better than `reflect.Indirect`

func KIndirectValueOf(data reflect.Value) reflect.Value {

  // safety
  if data.IsValid() {

    switch data.Kind() {

    case reflect.Interface: // Catch Any Type

      // unsafe, make infinity loop
      // recursion if contain interface, or ptr again
      return KIndirectValueOf(data.Elem())

    case reflect.Ptr:

      return data.Elem()
    }
  }

  return data
}
