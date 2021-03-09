package fasthttp

import (
	pub "github.com/edunx/rock-public-go"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"net"
)


func (self *Server) Start() error {

	pub.Out.Err("fasthttp start info %#v" , self.C)

	ln , err := self.Listen()
	if err != nil {
		return err
	}

	//init Server
	s := &fasthttp.Server{
		Handler: handler,
		TCPKeepalive: self.Keepalive(),
	}

	//注册虚拟路径
	cvr.Path(self.C.vhost)
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
