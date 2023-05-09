package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"leafy/app/models"
	"math"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"skfw/papaya/bunny/swag/method"
	"skfw/papaya/koala"
	m "skfw/papaya/koala/mapping"
	"skfw/papaya/koala/pp"
	"skfw/papaya/pigeon"
	"skfw/papaya/pigeon/drivers/common"
	"skfw/papaya/pigeon/drivers/postgresql"
	bacx "skfw/papaya/pigeon/templates/basicAuth/util"
	"strconv"
	"time"
)

type Printable func(args ...any)
type Testable func(t *Testing)

type Testing struct {
	Error Printable
	Log   Printable
}

type Case struct {
	Name    string
	Handler Testable
}

var userId, adminId string

func TestMergeData(t *Testing) {

	var err error

	os.Setenv("DB_USERNAME", "user")
	os.Setenv("DB_PASSWORD", "1234")
	os.Setenv("DB_NAME", "leafy")

	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_CHARSET", "utf8")
	os.Setenv("DB_TIMEZONE", "UTC")
	os.Setenv("DB_SECURE", "false")

	var conn common.DBConnectionImpl

	if conn, err = postgresql.DBConnectionNew(pigeon.InitLoadEnviron); err != nil {

		t.Error(err)
	}

	db := conn.GORM()

	db.AutoMigrate(&models.Users{}, &models.Sessions{}, &models.Cart{}, &models.Products{}, &models.Transactions{})

	users := db.Model(&models.Users{})
	carts := db.Model(&models.Cart{})
	sessions := db.Model(&models.Sessions{})
	products := db.Model(&models.Products{})
	category := db.Model(&models.Category{})
	categories := db.Model(&models.Categories{})
	nutrients := db.Model(&models.Nutrients{})
	transactions := db.Model(&models.Transactions{})

	var prepared *sql.DB

	prepared, _ = users.DB()
	prepared.Query("DELETE FROM users")

	prepared, _ = sessions.DB()
	prepared.Query("DELETE FROM sessions")

	prepared, _ = products.DB()
	prepared.Query("DELETE FROM products")

	prepared, _ = category.DB()
	prepared.Query("DELETE FROM category")

	prepared, _ = categories.DB()
	prepared.Query("DELETE FROM categories")

	prepared, _ = nutrients.DB()
	prepared.Query("DELETE FROM nutrients")

	prepared, _ = carts.DB()
	prepared.Query("DELETE FROM carts")

	prepared, _ = transactions.DB()
	prepared.Query("DELETE FROM transactions")

	pass, _ := bacx.HashPassword("User@1234")

	adminId = bacx.Idx(uuid.New())

	users.Create(map[string]any{
		"id":       adminId,
		"username": "admin",
		"email":    "admin@mail.co",
		"password": pass,
		"admin":    true,
	})

	userId = bacx.Idx(uuid.New())

	users.Create(map[string]any{
		"id":       userId,
		"username": "user",
		"email":    "user@mail.co",
		"password": pass,
		"admin":    true,
	})
}

type Map map[string]any
type Handler func(req *http.Request) error
type Maps []Map

func Mapping(data Map) m.KMapImpl {

	mm := m.KMap(data)
	return &mm
}

func ReqBody(method string, URL url.URL, payload Map, handler Handler) ([]byte, error) {

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

func Req(method string, URL url.URL, payload Map, handler Handler) (*Map, error) {

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

func ReqAr(method string, URL url.URL, payload Map, handler Handler) (Maps, error) {

	var err error
	var buff []byte

	if buff, err = ReqBody(method, URL, payload, handler); err != nil {

		return nil, err
	}

	var data Maps
	data = make(Maps, 0)

	if err = json.Unmarshal(buff, &data); err != nil {

		return nil, err
	}

	return data, nil
}

func Token(token string) Handler {

	return func(req *http.Request) error {

		req.Header.Set("Authorization", "Bearer "+token)

		return nil
	}
}

var origin = url.URL{
	Scheme: "http",
	Host:   "localhost:8000",
}

var token, tokenAdmin string

func ValueToInt(value any) int {

	val := pp.KIndirectValueOf(value)

	if val.IsValid() {

		ty := val.Type()

		switch ty.Kind() {

		case reflect.Float64:

			return int(m.KValueToFloat(value))
		}

		return int(m.KValueToInt(value))
	}

	return 0
}

func TestLogin(t *Testing) {

	var err error
	var res *Map

	origin.Path = "/api/v1/users/login"

	if res, err = Req(method.POST, origin, Map{
		"username": "user",
		"email":    "*",
		"password": "User@1234",
	}, nil); err != nil {

		t.Error(err)
	}

	token = m.KValueToString(Mapping(*res).Get("token"))

	if token == "" {

		t.Error("token empty")
	}

	t.Log("token_admin", token)
}

func TestAdminLogin(t *Testing) {

	var err error
	var res *Map

	origin.Path = "/api/v1/users/login"

	if res, err = Req(method.POST, origin, Map{
		"username": "admin",
		"email":    "*",
		"password": "User@1234",
	}, nil); err != nil {

		t.Error(err)
	}

	tokenAdmin = m.KValueToString(Mapping(*res).Get("token"))

	if tokenAdmin == "" {

		t.Error("token empty")
	}

	t.Log("token_admin", tokenAdmin)
}

func TestAdminAddProduct1(t *Testing) {

	var err error
	var res *Map

	origin.Path = "/api/v1/admin/product"

	if res, err = Req(method.POST, origin, Map{
		"name":        "Pineapple",
		"description": "A tropical plant with an edible fruit and the most economically significant plant in the family Bromeliaceae.",
		"price":       12,
		"stocks":      128,
		"categories":  []string{"pineapple", "raw"},
	}, Token(tokenAdmin)); err != nil {

		t.Error(err)
	}

	if check := m.KValueToBool(Mapping(*res).Get("error")); check {

		t.Error("failed create product")
	}

	t.Log("create product")
}

func TestAdminAddProduct2(t *Testing) {

	var err error
	var res *Map

	origin.Path = "/api/v1/admin/product"

	if res, err = Req(method.POST, origin, Map{
		"name":        "Apple",
		"description": "An edible fruit produced by an apple tree (Malus Domestica).",
		"price":       6,
		"stocks":      240,
		"categories":  []string{"apple", "raw"},
	}, Token(tokenAdmin)); err != nil {

		t.Error(err)
	}

	if check := m.KValueToBool(Mapping(*res).Get("error")); check {

		t.Error("failed create product")
	}

	t.Log("create product")
}

func TestAdminAddProduct3(t *Testing) {

	var err error
	var res *Map

	origin.Path = "/api/v1/admin/product"

	if res, err = Req(method.POST, origin, Map{
		"name":        "Papaya",
		"description": "The papaya, papaw, or pawpaw is the plant species Carica papaya, one of the 21 accepted species in the genus Carica of the family Caricaceae.",
		"price":       23,
		"stocks":      64,
		"categories":  []string{"papaya", "raw"},
	}, Token(tokenAdmin)); err != nil {

		t.Error(err)
	}

	if check := m.KValueToBool(Mapping(*res).Get("error")); check {

		t.Error("failed create product")
	}

	t.Log("create product")
}

func TestAdminAddProduct4(t *Testing) {

	var err error
	var res *Map

	origin.Path = "/api/v1/admin/product"

	if res, err = Req(method.POST, origin, Map{
		"name":        "Orange",
		"description": "An orange is a fruit of various citrus species in the family Rutaceae, it primarily refers to Citrus × sinensis, which is also called sweet orange.",
		"price":       6,
		"stocks":      120,
		"categories":  []string{"orange", "raw"},
	}, Token(tokenAdmin)); err != nil {

		t.Error(err)
	}

	if check := m.KValueToBool(Mapping(*res).Get("error")); check {

		t.Error("failed create product")
	}

	t.Log("create product")
}

func TestProductCatchAll(t *Testing) {

	var err error
	var res Maps

	origin.Path = "/api/v1/products"
	origin.RawQuery = "page=1&size=10"

	if res, err = ReqAr(method.GET, origin, nil, nil); err != nil {

		t.Error(err)
	}

	if len(res) == 0 {

		t.Error("data is empty")
	}

	t.Log("data", res)
}

func TestProductCart(t *Testing) {

	var err error
	var res Maps

	origin.Path = "/api/v1/products"
	origin.RawQuery = "page=1&size=10"

	if res, err = ReqAr(method.GET, origin, nil, nil); err != nil {

		t.Error(err)
	}

	if len(res) == 0 {

		t.Error("data is empty")
	}

	origin.Path = "/api/v1/users/cart"

	var check *Map
	var found bool

	for _, product := range res {

		found = false

		productId := m.KValueToString(Mapping(product).Get("id"))
		productName := m.KValueToString(Mapping(product).Get("name"))

		switch productName {
		case "Papaya", "Apple":

			if check, err = Req(method.POST, origin, Map{
				"productId": productId,
				"qty":       12,
			}, Token(token)); err != nil {

				t.Error(err)
				break
			}

			found = true

			break

		case "Pineapple":

			if check, err = Req(method.POST, origin, Map{
				"productId": productId,
				"qty":       3,
			}, Token(token)); err != nil {

				t.Error(err)
				break
			}

			found = true

			break
		}

		if found {

			t.Log(check)

			if m.KValueToBool(Mapping(*check).Get("error")) {

				t.Error("failed add cart")
				break
			}
		}
	}

	t.Log("data", res)
}

func TestCartCatchAll(t *Testing) {

	var err error
	var res Maps

	origin.Path = "/api/v1/users/carts"
	origin.RawQuery = "page=1&size=10"

	if res, err = ReqAr(method.GET, origin, nil, Token(token)); err != nil {

		t.Error(err)
	}

	if len(res) == 0 {

		t.Error("data is empty")
	}

	t.Log(res)
}

func TestBill(t *Testing) {

	var err error
	var res *Map

	origin.Path = "/api/v1/users/bill"

	if res, err = Req(method.GET, origin, nil, Token(token)); err != nil {

		t.Error(err)
	}

	bill := ValueToInt(Mapping(*res).Get("pay"))

	if bill != 384 {

		t.Error("bill incorrect value")
	}

	t.Log(res)
}

func TestTopUp(t *Testing) {

	var err error
	var res *Map

	origin.Path = "/api/v1/admin/topup"
	origin.RawQuery = "userId=" + userId

	if res, err = Req(method.POST, origin, Map{

		"balance": 1000,
	}, Token(token)); err != nil {

		t.Error(err)
	}

	t.Log(res)

	if check := m.KValueToBool(Mapping(*res).Get("error")); check {

		t.Error("failed topup")
	}

	t.Log("topup")
}

func TestPay(t *Testing) {

	var err error
	var res *Map

	origin.Path = "/api/v1/users/transaction"

	if res, err = Req(method.POST, origin, Map{

		"payment_method": "visa-credit-card",
	}, Token(token)); err != nil {

		t.Error(err)
	}

	t.Log(res)

	if check := m.KValueToBool(Mapping(*res).Get("error")); check {

		t.Error("failed create transaction")
	}

	t.Log("create transaction")
}

func TestPayCatchAll(t *Testing) {

	var err error
	var res Maps

	origin.Path = "/api/v1/users/transactions"
	origin.RawQuery = "page=1&size=10"

	if res, err = ReqAr(method.GET, origin, nil, Token(token)); err != nil {

		t.Error(err)
	}

	if len(res) == 0 {

		t.Error("data is empty")
	}

	t.Log(res)
}

var cases = []Case{

	{"Test Merge Data", TestMergeData},
	{"Test Login", TestLogin},
	{"Test Admin Login", TestAdminLogin},
	{"Test Admin Add Product 1", TestAdminAddProduct1},
	{"Test Admin Add Product 2", TestAdminAddProduct2},
	{"Test Admin Add Product 3", TestAdminAddProduct3},
	{"Test Admin Add Product 4", TestAdminAddProduct4},
	{"Test Product Catch All", TestProductCatchAll},
	{"Test Product Cart", TestProductCart},
	{"Test Cart Catch All", TestCartCatchAll},
	{"Test Bill", TestBill},
	{"Test TopUp", TestTopUp},
	{"Test Pay", TestPay},
	{"Test Pay Catch All", TestPayCatchAll},
}

func main() {

	var failed bool

	console := koala.KConsoleNew()

	failed = false
	testing := &Testing{
		Error: func(args ...any) {
			console.Error(args...)
			failed = true
		},
		Log: console.Log,
	}

	var k, n int
	var cas Case

	k, n = 1, len(cases)

	for _, cas = range cases {

		console.Log(
			console.Text("TEST", koala.ColorYellow, koala.ColorBlack, koala.StyleBold),
			console.Text(cas.Name, koala.ColorGreen, koala.ColorBlack, koala.StyleBold),
		)

		cas.Handler(testing)

		if failed {

			break
		}

		k++
	}

	percent := int(math.Min(float64(k)/float64(n)*100, 100))

	console.Log(console.Text("PASS", koala.ColorCyan, koala.ColorBlack, koala.StyleBold), console.Text(strconv.Itoa(percent)+"%", koala.ColorYellow, koala.ColorBlack, koala.StyleBold))
}
