package gen

import (
	"PapayaNet/papaya/koala"
	"reflect"
	"strconv"
	"strings"
)

type KMapIteration struct {
	KIterationImpl[string, any]
}

type KMapIterationNextHandler KIterationNextHandler[string, any]

type KMapIterationImpl interface {
	KIterationImpl[string, any]
}

func KMapStopIteration() KMapIterationImpl {

	return &KMapIteration{
		&KIteration[string, any]{
			stopIter: true,
		},
	}
}

func KMapIterable(mapping any) KMapIterationImpl {

	val := koala.KIndirectValueOf(reflect.ValueOf(mapping))
	ty := val.Type()

	if val.IsValid() {

		switch val.Kind() {
		case reflect.Array, reflect.Slice:

			n := val.Len()
			k := 0

			if n > 0 {

				return &KMapIteration{
					&KIteration[string, any]{
						NextHandler: func(v KIterationImpl[string, any]) error {

							if k < n {

								key := strconv.Itoa(k)
								value := val.Index(k).Interface()

								v.SetValues(key, value)
							}

							k += 1

							return nil
						},
					},
				}
			}

			break
		case reflect.Map:

			if ty.Key().Kind() == reflect.String {

				mapRange := val.MapRange()

				return &KMapIteration{
					&KIteration[string, any]{
						NextHandler: func(v KIterationImpl[string, any]) error {

							if hasNext := mapRange.Next(); hasNext {

								key := mapRange.Key().String()
								value := mapRange.Value().Interface()

								v.SetValues(key, value)
							}

							return nil
						},
					},
				}
			}

			break

		case reflect.Struct:

			n := ty.NumField()
			k := 0

			if n > 0 {

				return &KMapIteration{
					&KIteration[string, any]{
						NextHandler: func(v KIterationImpl[string, any]) error {

							if k < n {

								key := ty.Field(k).Name
								value := val.Field(k).Interface()

								v.SetValues(key, value)
							}

							k += 1

							return nil
						},
					},
				}
			}
		}
	}

	return KMapStopIteration()
}

type KMapPageIteration struct {
	KPageIterationImpl[string, any]
}

type KMapPageIterationImpl interface {
	KPageIterationImpl[string, any]
}

func KMapIterableDeeper(mapping any) KMapIterationImpl {

	//val := koala.KIndirectValueOf(reflect.ValueOf(mapping))
	//ty := val.Type()

	mapPageIteration := &KMapPageIteration{
		&KPageIteration[string, any]{},
	}
	mapPageIteration.Init()

	iter := KMapIterable(mapping)
	mapPageIteration.Add(iter)

	return &KMapIteration{
		&KIteration[string, any]{
			NextHandler: func(v KIterationImpl[string, any]) error {

				next := mapPageIteration.Wait()

				if next.HasNext() {

					enum := next.Enum()
					value := enum.Value()
					keys := mapPageIteration.Keys()

					if KMapHunt(value) {

						iter := KMapIterable(value)
						mapPageIteration.Add(iter)
					}

					v.SetValues(strings.Join(keys, "."), value)
				}

				return nil
			},
		},
	}
}
