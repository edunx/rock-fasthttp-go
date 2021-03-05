package fasthttp

import (
	"context"
	"github.com/edunx/lua"
	"sync"
)

var (
	cvr = vRCache{ pool: sync.Map{} }
	cvm = vMCache{ pool: sync.Map{} }
)

type (
	region interface {
		Search( string ) (int64 , []byte , error)
	}

	Config struct {
		listen    string
		vhost     string
		handler   string
		keepalive string
		protocol  string
		reuseport string
		unknown   string
		daemon    string
	}

	Server struct {
		C         Config
		region    region
	}
)

type (

	KeyVal struct {
		Key string
		Val string
	}

	thread struct {
		co  *lua.LState
		cancelFunc context.CancelFunc
	}

	vHandler struct {
		rule          []string
		header        []*KeyVal
		tag           string
		body          string
		code          int
		bodyEncode    string
		bodyEncodeMin int
		eof           string
		hook          *lua.LFunction

		Pool          *sync.Pool //worker pool
	}
)

type (
	vRouter struct {
		L       *lua.LState
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
	}
)


