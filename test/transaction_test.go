package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"io"
	"leafy/app/models"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"skfw/papaya/bunny/swag/method"
	m "skfw/papaya/koala/mapping"
	"skfw/papaya/koala/pp"
	"skfw/papaya/pigeon"
	"skfw/papaya/pigeon/drivers/common"
	"skfw/papaya/pigeon/drivers/postgresql"
	"skfw/papaya/pigeon/templates/basicAuth/repository"
	"testing"
	"time"
)

var userId, adminId string

func TestMergeData(t *testing.T) {

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
	transactions := db.Model(&models.Transactions{})

	var prepared *sql.DB

	prepared, _ = users.DB()
	prepared.Query("DELETE FROM users")

	prepared, _ = sessions.DB()
	prepared.Query("DELETE FROM sessions")

	prepared, _ = products.DB()
	prepared.Query("DELETE FROM products")

	prepared, _ = carts.DB()
	prepared.Query("DELETE FROM carts")

	prepared, _ = transactions.DB()
	prepared.Query("DELETE FROM transactions")

	pass, _ := repository.HashPassword("User@1234")

	adminId = repository.Idx(uuid.New())

	users.Create(map[string]any{
		"id":       adminId,
		"username": "admin",
		"email":    "admin@mail.co",
		"password": pass,
		"admin":    true,
	})

	userId = repository.Idx(uuid.New())

	users.Create(map[string]any{
		"id":       userId,
		"username": "user",
		"email":    "user@mail.co",
		"password": pass,
		"admin":    true,
	})

	products.Create(&models.Products{
		ID:          repository.Idx(uuid.New()),
		Name:        "Apple",
		Description: "Red Apple",
		Stocks:      100,
		Price:       decimal.NewFromInt(12),
	})

	products.Create(&models.Products{
		ID:          repository.Idx(uuid.New()),
		Name:        "Papaya",
		Description: "Papaya Mountain",
		Stocks:      120,
		Price:       decimal.NewFromInt(23),
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

func TestProductCatchAll(t *testing.T) {

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

func TestLogin(t *testing.T) {

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

func TestAdminLogin(t *testing.T) {

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

func TestAdminAddProduct(t *testing.T) {

	var err error
	var res *Map

	origin.Path = "/api/v1/admin/product"

	if res, err = Req(method.POST, origin, Map{
		"name":        "Pineapple",
		"description": "Pineapple Yellow from jungle",
		"price":       12,
		"stocks":      128,
	}, Token(tokenAdmin)); err != nil {

		t.Error(err)
	}

	if check := m.KValueToBool(Mapping(*res).Get("error")); check {

		t.Error("failed create product")
	}

	t.Log("create product")
}

func TestProductCart(t *testing.T) {

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

func TestCartCatchAll(t *testing.T) {

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

func TestBill(t *testing.T) {

	var err error
	var res *Map

	origin.Path = "/api/v1/users/bill"

	if res, err = Req(method.GET, origin, nil, Token(token)); err != nil {

		t.Error(err)
	}

	bill := ValueToInt(Mapping(*res).Get("pay"))

	if bill != 456 {

		t.Error("bill incorrect value")
	}

	t.Log(res)
}

func TestTopUp(t *testing.T) {

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

func TestPay(t *testing.T) {

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

func TestPayCatchAll(t *testing.T) {

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
