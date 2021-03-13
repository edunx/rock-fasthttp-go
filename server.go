package fasthttp

import (
	pub "github.com/edunx/rock-public-go"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"net"
)

func (self *Server) handler( ctx *fasthttp.RequestCtx ) {
	ctx.Logger().Printf("logger")

	vrr := cvr.load( pub.B2S( ctx.Host() ) )
	if vrr == nil {
		ctx.Response.SetStatusCode(500)
		ctx.Response.SetBody(pub.S2B("not found router"))
		return
	}

	r , ok := vrr.L.GetExdata().(*router.Router)
	if !ok {
		ctx.Response.SetStatusCode(500)
		ctx.Response.SetBody(pub.S2B("expect invalid router"))
		return
	}

	ctx.SetUserValue( "vrr" , vrr)
	r.Handler(ctx)
}

func (self *Server) Start() error {
	pub.Out.Err("fasthttp start info %#v" , self.C)

	ln , err := self.Listen()
	if err != nil {
		return err
	}

	//init Server
	s := &fasthttp.Server{
		Handler: self.handler,
		TCPKeepalive: self.Keepalive(),
	}

	//注册虚拟路径
	cvr.Path(self.C.routers)
	cvr.Unknown(self.C.unknown)

	cvm.Path( self.C.handler )
	go cvr.sync()
	go cvm.sync()

	if self.C.daemon == "on" {
		go s.Serve(ln)
	} else {
		s.Serve(ln)
	}

	return nil
}

func(self *Server) Keepalive() bool {
	if self.C.keepalive == "on" {
		return true
	}
	return false
}

func(self *Server) Listen() (net.Listener , error) {
	if self.C.reuseport == "on" {
		return reuseport.Listen(self.C.protocol , self.C.listen)
	}
	return net.Listen(self.C.protocol , self.C.listen)
}
