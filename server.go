package fasthttp

import (
	"github.com/edunx/lua"
	pub "github.com/edunx/rock-public-go"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"net"
)

func (self *Server) newLogger() {

	self.vlog = vlogger{}

	self.vlog.New(self.C.accessLog)
	switch self.C.accessFormat {
	case "json":
		self.vlog.encode = self.vlog.Json
	case "raw":
		self.vlog.encode = self.vlog.Raw
	default:
		self.vlog.encode = self.vlog.Raw
	}
}

func (self *Server) Logger( ctx *fasthttp.RequestCtx  , vrr *vRouter ) {

	off , ok := ctx.UserValue("access_push_off").(string)
	if ok && off == "off" {
		return
	}

	off , ok = vrr.L.ExData.Get("access_push_off").(string)
	if ok && off == "off" {
		return
	}

	if self.access == nil {
		return
	}

	vlog , ok := vrr.L.ExData.Get("logger").(*vlogger)
	if ok {
		self.access.Push( vlog.encode(ctx) )
		//pub.Out.Err("%s" , vlog.encode(ctx))
		return
	}

	self.access.Push(self.vlog.encode(ctx))
	//pub.Out.Err("%s" , self.vlog.encode(ctx))
}

func (self *Server) Region( L *lua.LState , ctx *fasthttp.RequestCtx ) {
	var addr string
	var ok bool
	var key string

	key , ok = ctx.UserValue("region").(string)
	if ok {
		goto RUN
	}

	key , ok = L.ExData.Get("region").(string)
	if ok {
		goto RUN
	}

	key = self.C.accessRegion
RUN:

	switch key {
	case "off":
		return
	case "remote_addr":
		addr = ctx.RemoteIP().String()
	default:
		addr = string(ctx.Request.Header.Peek( key ))
	}

	if addr == "" {
		return
	}

	cityid , info , err := self.region.Search(addr)

	if err != nil {
		return
	}

	ctx.SetUserValue("region_city" , cityid)
	ctx.SetUserValue("region_info" ,  info)
}

func (self *Server) handler( ctx *fasthttp.RequestCtx ) {

	vrr := cvr.load( pub.B2S( ctx.Host() ) )
	if vrr == nil {
		ctx.Response.SetStatusCode(500)
		ctx.Response.SetBody(pub.S2B("not found router"))
		return
	}


	//获取解析的lua虚拟机
	r  := CheckExDataRouter( vrr.L )

	//获取运行lua 虚拟机环境
	ctx.SetUserValue("vctx" , &vContext{vrr:vrr})

	//处理业务逻辑,匹配路由模式
	r.Handler(ctx)

	//解析用户地址位置
	self.Region(vrr.L , ctx)

	//操作日志
	self.Logger(ctx , vrr)

}

func (self *Server) Start() error {
	pub.Out.Err("fasthttp start info %#v" , self.C)

	ln , err := self.Listen()
	if err != nil {
		return err
	}

	self.FServer = &fasthttp.Server{
		Handler: self.handler,
		TCPKeepalive: self.Keepalive(),
	}

	//注册虚拟路径
	cvr.Path(self.C.routers)
	cvr.Unknown(self.C.unknown)

	//注册日志模块
	self.newLogger()

	cvm.Path( self.C.handler )
	go cvr.sync()
	go cvm.sync()

	if self.C.daemon == "on" {
		go self.FServer.Serve(ln)
	} else {
		self.FServer.Serve(ln)
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
