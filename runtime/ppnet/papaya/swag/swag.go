package swag

import (
	"PapayaNet/papaya/koala"
	"PapayaNet/papaya/koala/kornet"
)

const (
	VersionMajor = 2
	VersionMinor = 0
	VersionPatch = 0
)

type PnSwagParameter struct {
	Name        string
	In          string
	Description string
	Required    bool
	Type        string
}

// if you use `array`, replace `structure` with `items`
// if you use `object`, replace `structure` with `properties`

type PnSwagValue struct {
	Name      string
	Type      string
	Structure *PnSwagValue // if you use `array` or `object` type
}

type PnSwagSchema struct {
	Type      string
	Structure *PnSwagValue
}

type PnSwagResponse struct {
	Status      int
	Description string
	Schema      *PnSwagSchema
}

type PnSwagPathMethodDescription struct {
	Tags        []string
	Summary     string
	Description string
	OperationId string // TODO: lookup name function as id
	Consumes    []string
	Produces    []string
	Parameters  []PnSwagParameter
	Responses   []PnSwagResponse
	Deprecated  bool
}

type PnSwagPathMethod struct {
	Method      string
	Description *PnSwagPathMethodDescription
}

type PnSwagPath struct {
	Path   string
	Method []PnSwagPathMethod
}

type PnSwagTermOfService struct {
	URL string // no prop
}

type PnSwagContact struct {
	Email string
}

type PnSwagLicence struct {
	Name string
	URL  string
}

type PnSwagInfo struct {
	Description   string
	Version       string
	Title         string
	TermOfService *PnSwagTermOfService // no prop
	Contact       *PnSwagContact
	License       *PnSwagLicence
}

type PnSwag struct {
	Version *koala.KVersion
	Schemes []string
	Info    *PnSwagInfo
	Path    []PnSwagPath
}

type PnSwagRespHandler func(ctx *kornet.KResponse) error

type PnSwagRouterImpl interface {
	Get(path string, params koala.KMap, body koala.KMap, resp koala.KMap, handler PnSwagRespHandler)
	Head(path string, params koala.KMap, body koala.KMap, resp koala.KMap, handler PnSwagRespHandler)
	Post(path string, params koala.KMap, body koala.KMap, resp koala.KMap, handler PnSwagRespHandler)
	Put(path string, params koala.KMap, body koala.KMap, resp koala.KMap, handler PnSwagRespHandler)
	Delete(path string, params koala.KMap, body koala.KMap, resp koala.KMap, handler PnSwagRespHandler)
	Connect(path string, params koala.KMap, body koala.KMap, resp koala.KMap, handler PnSwagRespHandler)
	Options(path string, params koala.KMap, body koala.KMap, resp koala.KMap, handler PnSwagRespHandler)
	Trace(path string, params koala.KMap, body koala.KMap, resp koala.KMap, handler PnSwagRespHandler)
}

type PnSwagGroupImpl interface {
	Router() PnSwagRouterImpl
}

type PnSwagImpl interface {
	Init() error
	Group(name string) PnSwagGroupImpl
	Router() PnSwagRouterImpl
}

func (swag *PnSwag) Init() {

	swag.Version = koala.KVersionNew(
		VersionMajor,
		VersionMinor,
		VersionPatch,
	)
}
