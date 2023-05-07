package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	m "skfw/papaya/koala/mapping"
	"skfw/papaya/koala/pp"
	"strings"
	"time"
)

type Map map[string]any
type Handler func(req *http.Request) error
type Maps []Map

// Nutrient API

func NAPI() *url.URL {

	return &url.URL{
		Scheme:   "https",
		Host:     "api.nal.usda.gov",
		Path:     "fdc/v1/foods/search",
		RawQuery: "dataType=Survey%20%28FNDDS%29&pageSize=1&pageNumber=1&sortBy=dataType.keyword&sortOrder=asc&API_KEY=" + strings.Trim(os.Getenv("NUTRIENT_API_KEY"), " "),
	}
}

func ReqBody(method string, URL *url.URL, payload Map, handler Handler) ([]byte, error) {

	var err error
	var req *http.Request
	var res *http.Response

	body := new(bytes.Buffer)
	encoder := json.NewEncoder(body)

	if err = encoder.Encode(payload); err != nil {

		return nil, err
	}

	if req, err = http.NewRequest(method, URL.String(), body); err != nil {

		return nil, err
	}

	if handler != nil {

		if err = handler(req); err != nil {

			return nil, err
		}
	}

	client := http.Client{Timeout: time.Second * 15}

	if res, err = client.Do(req); err != nil {

		return nil, err
	}

	defer func(Body io.ReadCloser) {

		if err = Body.Close(); err != nil {

			panic(err)
		}
	}(res.Body)

	var data []byte

	if data, err = io.ReadAll(res.Body); err != nil {

		return nil, err
	}

	return data, nil
}

func Req(method string, URL *url.URL, payload Map, handler Handler) (*Map, error) {

	var err error
	var buff []byte

	if buff, err = ReqBody(method, URL, payload, handler); err != nil {

		return nil, err
	}

	var data Map

	if err = json.Unmarshal(buff, &data); err != nil {

		return nil, err
	}

	return &data, nil
}

type Nutrient struct {
	Name  string `json:"name"`
	Unit  string `json:"unit"`
	Value int    `json:"value"`
}

func NutrientAPI(categories []string) []Nutrient {

	var err error
	var data *Map

	temp := make([]Nutrient, 0)

	URL := NAPI()
	URL.RawQuery += "&query=" + url.QueryEscape(strings.Join(categories, " "))

	if data, err = Req("GET", URL, nil, nil); err != nil {

		return temp
	}

	mm := m.KMap(*data)

	val := pp.KIndirectValueOf(mm.Get("foods.0.foodNutrients"))

	nutrientNameOf := reflect.ValueOf("nutrientName")
	unitNameOf := reflect.ValueOf("unitName")
	valueOf := reflect.ValueOf("value")

	if val.IsValid() {

		ty := val.Type()

		switch ty.Kind() {
		case reflect.Array, reflect.Slice:

			for i := 0; i < val.Len(); i++ {

				elem := pp.KIndirectValueOf(val.Index(i))

				if elem.IsValid() {

					tyElem := elem.Type()

					switch elem.Kind() {
					case reflect.Map:

						if tyElem.Key().Kind() == reflect.String {

							nutrientName := m.KValueToString(elem.MapIndex(nutrientNameOf))
							unitName := m.KValueToString(elem.MapIndex(unitNameOf))
							value := ValueToInt(elem.MapIndex(valueOf))

							if value > 0 {

								temp = append(temp, Nutrient{nutrientName, unitName, value})
							}
						}

						break
					}
				}
			}

			break
		}
	}

	return temp
}
