package mapping

import (
	"PapayaNet/papaya/koala"
	"reflect"
	"strconv"
	"strings"
)

func KMapKeys(mapping any) []string {

	tokens := make([]string, 0)

	value := koala.KIndirectValueOf(reflect.ValueOf(mapping))
	ty := value.Type()

	if value.IsValid() {

		switch value.Kind() {
		case reflect.Array, reflect.Slice:

			for i := 0; i < value.Len(); i++ {

				tokens = append(tokens, strconv.Itoa(i))
			}

			break

		case reflect.Map:

			if ty.Key().Kind() == reflect.String {

				mapIter := value.MapRange()

				for mapIter.Next() {

					key := mapIter.Key()

					tokens = append(tokens, key.String())
				}
			}

			break
		case reflect.Struct:

			// TODO: not implemented yet

			break
		}
	}

	return tokens
}

func KMapKeysDeeper(mapping any) []string {

	tokens := make([]string, 0)

	value := koala.KIndirectValueOf(reflect.ValueOf(mapping))
	ty := value.Type()

	if value.IsValid() {

		switch value.Kind() {
		case reflect.Array, reflect.Slice:

			for i := 0; i < value.Len(); i++ {

				tokens = append(tokens, strconv.Itoa(i))
			}

			break

		case reflect.Map:

			if ty.Key().Kind() == reflect.String {

				mapIter := value.MapRange()

				for mapIter.Next() {

					key := mapIter.Key()
					value := mapIter.Value()

					fk := key.String()
					keys := KMapKeysDeeper(value.Interface())

					if len(keys) > 0 {

						for _, k := range keys {

							tokens = append(tokens, fk+"."+k)
						}

						continue
					}

					tokens = append(tokens, fk)
				}
			}

			break
		case reflect.Struct:

			// TODO: not implemented yet

			n := value.NumField()

			for i := 0; i < n; i++ {

				//f := value.Field(i)
				ft := ty.Field(i)

				//ft.Name, f.Type(), ft.Tag

				tokens = append(tokens, ft.Name)
			}

			break
		}
	}

	return tokens
}

func KMapValues(mapping any) []any {

	values := make([]any, 0)

	// value is KMap (fixed)
	value := koala.KIndirectValueOf(reflect.ValueOf(mapping))
	ty := value.Type()

	if value.IsValid() {

		switch value.Kind() {
		case reflect.Map:

			if ty.Key().Kind() == reflect.String {

				mapIter := value.MapRange()

				for mapIter.Next() {

					value := koala.KIndirectValueOf(mapIter.Value())
					values = append(values, value)
				}
			}

			break
		case reflect.Struct:

			// TODO: not implemented yet

			n := value.NumField()

			for i := 0; i < n; i++ {

				f := value.Field(i)
				//ft := ty.Field(i)
				//ft.Name (FieldName), f.Interface() (Value), ft.Tag (FieldTag)

				values = append(values, f.Interface())
			}

			break
		}
	}

	return values
}

func KMapGetValue(name string, mapping any) any {

	tokens := strings.Split(name, ".")

	var value any
	value = mapping

	i := 0
	n := len(tokens)

	if n == 0 {

		return nil
	}

	for {

		if value == nil {

			return nil
		}

		if n <= i {
			break
		}

		token := tokens[i]

		val := koala.KIndirectValueOf(reflect.ValueOf(value))
		ty := val.Type()

		if val.IsValid() {

			switch ty.Kind() {

			case reflect.Array, reflect.Slice:

				if index, err := strconv.Atoi(token); err == nil {

					if index < val.Len() {

						v := val.Index(index)
						value = v.Interface()
						break
					}
				}

				value = nil

				break

			case reflect.Map:

				if ty.Key().Kind() == reflect.String {

					m := false

					mapIter := val.MapRange()

					for mapIter.Next() {

						key := mapIter.Key().String()

						if token == key {

							//value := reflect.Indirect(mapIter.Value())
							v := mapIter.Value()
							value = v.Interface()
							m = true
							break
						}
					}

					if !m {

						value = nil
					}
				}

				break

			case reflect.Struct:

				// TODO: not implemented yet

				m := false

				for j := 0; j < val.NumField(); j++ {

					f := val.Field(j)
					ft := ty.Field(j)

					if token == ft.Name {

						value = f.Interface()
						m = true
						break
					}
				}

				if !m {

					value = nil
				}

				break
			}
		}

		i++
	}

	return value
}

func KMapSetValue(name string, data any, mapping koala.KMap) bool {

	tokens := strings.Split(name, ".")

	value := reflect.ValueOf(mapping)

	i := 0
	n := len(tokens)

	if n == 0 {

		return false
	}

	// set value as a pointer ? im don't thing so, because can set it without catch pointer
	// `array`, `slice`, or `map` that combine ptr on inside memory
	// that get previous and set value

	for {

		// that have problem, type hide on inside interface and ptr
		// how to solve, is catch all elems on inside interface and set back in var `value`
		// a reset type in var `ty`
		value = koala.KIndirectValueOf(value)

		if !value.IsValid() {

			return false
		}

		if n <= i {
			break
		}

		token := tokens[i]

		// update type of temporary
		ty := value.Type()

		if value.IsValid() {

			// lookup data on `map`
			switch value.Kind() {

			case reflect.Array, reflect.Slice:

				if index, err := strconv.Atoi(token); err == nil {

					if index < value.Len() {

						// get previous, set value
						if i+1 == n {

							// get previous, to set value on inside `array` or `slice`
							value.Index(index).Set(reflect.ValueOf(data))
							return true
						}

						v := value.Index(index)
						value = v
						break
					}

					// index out of bound
					return false
				}

				// NaN
				return false

			case reflect.Map:
				if ty.Key().Kind() == reflect.String {

					m := false

					mapIter := value.MapRange()

					for mapIter.Next() {

						key := mapIter.Key()

						if token == key.String() {

							// get previous, set value
							if i+1 == n {

								// get previous, to set value on inside `map`
								value.SetMapIndex(key, reflect.ValueOf(data))
								return true
							}

							v := mapIter.Value()
							value = v

							m = true
							break
						}
					}

					if !m {

						return false
					}

					break
				}

				// bad key
				return false

			case reflect.Struct:

				// TODO: not implemented yet

				n := value.NumField()

				for i := 0; i < n; i++ {

					f := value.Field(i)
					ft := ty.Field(i)

					if token == ft.Name {

						f.Set(reflect.ValueOf(data))
						return true
					}
				}

				return false
			}
		}

		i++
	}

	// other than `array`, `slice`, or `map`
	// can't set value as ptr
	return false
}

// TODO: fix slicing

func KMapDelValue(name string, mapping koala.KMap) bool {

	tokens := strings.Split(name, ".")

	var value, prev reflect.Value

	value = reflect.ValueOf(mapping)
	prev = reflect.Value{}

	var prevToken string

	i := 0
	n := len(tokens)

	if n == 0 {

		return false
	}

	// set value as a pointer ? im don't thing so, because can set it without catch pointer
	// `array`, `slice`, or `map` that combine ptr on inside memory
	// that get previous and set value

	for {

		if value.Interface() == nil {

			return false
		}

		if n <= i {
			break
		}

		token := tokens[i]

		// that have problem, type hide on inside interface and ptr
		// how to solve, is catch all elems on inside interface and set back in var `value`
		// a reset type in var `ty`
		value = koala.KIndirectValueOf(value)

		// get type
		ty := value.Type()

		if value.IsValid() {

			// lookup data on `map`
			switch value.Kind() {

			case reflect.Array, reflect.Slice:

				if index, err := strconv.Atoi(token); err == nil {

					if index < value.Len() {

						// get previous, set value
						if i+1 == n {

							s := value.Len()

							// get previous, to set value on inside `array` or `slice`
							L := value.Slice(0, index)
							R := value.Slice(index+1, s)

							k := L.Len() + R.Len()

							// ERROR: problem
							data := reflect.MakeSlice(value.Type(), k, k)

							// merging

							for j := 0; j < L.Len(); j++ {

								data.Index(j).Set(L.Index(j))
							}

							for j := 0; j < R.Len(); j++ {

								data.Index(j + L.Len()).Set(R.Index(j))
							}

							// end

							// save on previous value as ref
							if reflect.ValueOf(prevToken).IsValid() {

								switch prev.Kind() {

								case reflect.Array, reflect.Slice:

									if index, err := strconv.Atoi(prevToken); err == nil {

										prev.Index(index).Set(data)
										return true
									}

									return false

								case reflect.Map:

									prev.SetMapIndex(reflect.ValueOf(prevToken), data)
									break

								case reflect.Struct:

									prev.FieldByName(prevToken).Set(data)
									break
								}

								return true
							}

							// try save on current elem
							// panic: reflect: reflect.Value.Set using unaddressable value
							//value.Set(data)
							//if reflect.DeepEqual(data, value.Interface()) {
							//
							//  return true
							//}

							return false
						}

						v := value.Index(index)
						prev = value
						value = v

						break
					}

					// index out of bound
					return false
				}

				// NaN
				return false

			case reflect.Map:

				if ty.Key().Kind() == reflect.String {

					m := false

					mapIter := value.MapRange()

					for mapIter.Next() {

						key := mapIter.Key()

						if token == key.String() {

							// get previous, set value
							if i+1 == n {

								// get previous, to delete value on inside `map`
								value.SetMapIndex(key, reflect.Value{}) // that don't know how it works
								return true
							}

							v := mapIter.Value()
							prev = value
							value = v

							m = true
							break
						}
					}

					if !m {

						return false
					}

					break
				}

				// bad key
				return false

			case reflect.Struct:

				// TODO: not implemented yet

				n := value.NumField()

				for i := 0; i < n; i++ {

					f := value.Field(i)
					ft := ty.Field(i)

					if token == ft.Name {

						// make it zero value
						f.Set(reflect.Zero(f.Type()))
						return true
					}
				}

				return false
			}
		}

		prevToken = token
		i++
	}

	// other than `array`, `slice`, or `map`
	// can't set value as ptr
	return false
}
