package papaya

import (
	"PapayaNet/papaya/db"
	"PapayaNet/papaya/db/drivers/mysql"
	"PapayaNet/papaya/koala"
	"PapayaNet/papaya/koala/kornet"
	"crypto/tls"
	"errors"
	"github.com/labstack/echo/v4"
	gLog "github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"net"
	"net/http"
	"strconv"
)

type PapayaNet struct {
	DB            *gorm.DB
	Echo          *echo.Echo
	Version       koala.KVersionImpl
	Console       koala.KConsoleImpl
	HasInitialed  bool
	HasEchoClosed bool
	HasDBClosed   bool
	HasClosed     bool
}

type PapayaNetImpl interface {
	Init() error
	EnvLoader(environ PnEnvImpl) error
	GetConsole() koala.KConsoleImpl
	Serve(addr string, port uint16) error
	EchoClose() error
	DBClose() error
	Close() error
}

func (pn *PapayaNet) GetConsole() koala.KConsoleImpl {

	if pn.Console != nil {

		return pn.Console
	}

	return koala.KConsoleNew()
}

func (pn *PapayaNet) EchoClose() error {

	if e := pn.Echo; e != nil {

		err := e.Close()
		if err != nil {

			return err
		}

		pn.HasEchoClosed = true
		return nil
	}

	return nil
}

func (pn *PapayaNet) DBClose() error {

	if dbGorm := pn.DB; dbGorm != nil {

		d, e := dbGorm.DB()

		if e != nil {

			return e
		}

		err := d.Close()

		if err != nil {

			return err
		}

		pn.HasDBClosed = true
		return nil
	}

	return nil
}

func (pn *PapayaNet) Close() error {

	if !pn.HasEchoClosed {

		if err := pn.EchoClose(); err != nil {

			return err
		}

		pn.HasEchoClosed = true
	}

	if !pn.HasDBClosed {

		if err := pn.DBClose(); err != nil {

			return err
		}

		pn.HasDBClosed = true
	}

	pn.HasClosed = true

	return nil
}

func (pn *PapayaNet) Init() error {

	// skipping new initial
	if pn.HasInitialed {

		return nil
	}

	pn.HasInitialed = true

	if pn.Version == nil {

		pn.Version = koala.KVersionNew(
			VersionMajor,
			VersionMinor,
			VersionPatch)
	}

	banner := Banner(pn.Version)

	if pn.Console == nil {

		console := koala.KConsoleNew()
		//console.Colorful = true // force use colored
		//console.Silent = true // force silent mode
		console.Listen(func(info int) error {

			switch info {
			case koala.TypeError, koala.TypeWarn:
				console.Warn("Warning!!")
			}
			return nil
		})

		pn.Console = console
	}

	pn.Console.Log(
		pn.Console.EOL(),
		pn.Console.Text(
			banner,
			koala.ColorGreen,
			koala.ColorBlack,
			koala.StyleBold),
	)

	// load env_module
	if err := pn.EnvLoader(&PnDotEnv{}); err != nil {

		pn.Console.Error(err)
	}

	if pn.Echo == nil {

		pn.Echo = echo.New()
		//pn.Echo.Pre(middleware.Logger())
		//pn.Echo.Use(middleware.Recover())
		pn.Echo.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(ctx echo.Context) error {

				pn.Console.Main(func() error {

					request := ctx.Request()
					URL := kornet.KRequestGetURL(request)
					pn.Console.Log(URL.String()) // Logger
					return nil
				})

				return next(ctx)
			}
		})
		pn.Echo.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(ctx echo.Context) error {

				response := ctx.Response()
				response.Header().Set("Server", "PapayaNet "+pn.Version.String())

				return next(ctx)
			}
		})
	}

	if pn.DB == nil {

		client := &mysql.DBClient{}
		myClient, err := client.Init(db.InitLoadEnv)
		if err != nil {

			pn.Console.Error(err)
		}
		if myClient != nil {

			pn.DB = myClient
		}
	}

	pn.Echo.HideBanner = true
	if echoLog, ok := pn.Echo.Logger.(*gLog.Logger); ok {

		echoLog.SetHeader("[${time_rfc3339_nano}] [${level}]")
	}

	pn.Console.Error("testing ...")

	return nil
}

func (pn *PapayaNet) EnvLoader(environ PnEnvImpl) error {

	return environ.Load()
}

func (pn *PapayaNet) Serve(addr string, port uint16) error {

	if port < 1024 {

		return errors.New("params `port` prevent to use! that mean port use for core systems")
	}

	sPort := strconv.Itoa(int(port))

	pHttp := &http.Server{Addr: addr + ":" + sPort}
	pHttp.Handler = pn.Echo
	pHttp.ErrorLog = pn.Echo.StdLogger

	if pHttp.TLSConfig == nil {

		Listener, err := net.Listen("tcp", pHttp.Addr)

		if err != nil {

			return err
		}

		pn.Echo.Listener = Listener

		pn.Console.Log("Server started on " + addr + ":" + sPort)
		return pHttp.Serve(pn.Echo.Listener)
	}

	TLSListener, err := tls.Listen("tcp", pHttp.Addr, pHttp.TLSConfig)

	if err != nil {

		return err
	}

	pn.Echo.TLSListener = TLSListener

	pn.Console.Log("Server started on " + addr + ":" + sPort)
	return pHttp.Serve(pn.Echo.TLSListener)
}
