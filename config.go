package fasthttp

import (
	"context"
	"github.com/edunx/lua"
	base "github.com/edunx/rock-base-go"
	tp "github.com/edunx/rock-transport-go"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"sync"
)

const (
	VHANDLER  int = iota //vhander struct
	VHSTRING             //string handler
)

var (
	cvr = vRCache{ pool: sync.Map{} }
	cvm = vMCache{ pool: sync.Map{} }
)

type (
	region interface {
		Search( string ) (int64 , []byte , error)
	}

	vPush interface {
		push( interface{} )
	}

	Config struct {
		listen       string
		routers      string
		handler      string
		keepalive    string
		protocol     string
		reuseport    string
		unknown      string
		daemon       string
		accessLog    string
		accessFormat string
		accessRegion string
	}

	Server struct {
		C            Config
		region       region
		FServer      *fasthttp.Server
		vlog         vlogger
		access       tp.Tunnel
	}
)

type (

	thread struct {
		co  *lua.LState
		cancelFunc context.CancelFunc
	}

	vContext struct {
		vrr *vRouter
		vth *thread
	}

	vRule struct {
		name   string
		method string
		value  interface{}
	}

	vHandler  struct {
		count         int
		rule          []string
		header        []*base.KeyVal
		body          string
		code          int
		eof           string
		hook          *lua.LFunction
	}

	vHandlerChains struct {
		data  []interface{}
		mask  []int

		cap   int
	}

	vlogger struct {
		format      []string
		encode      func(*fasthttp.RequestCtx) []byte
	}

)

type (
	//转化
	Router	router.Router

	vRouter struct {
		L           *lua.LState
		Co          sync.Pool

		modTime int64
		name    string
	}

	vRCache struct {
		pool    sync.Map
		unknown *vRouter
		path    string
	}
)

type(
	vModule struct {
		handler *vHandler
		modTime int64
		name    string
	}

	vMCache struct {
		pool    sync.Map
		path    string
		L       *lua.LState
		once    sync.Once
	}
)



