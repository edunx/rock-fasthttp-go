package fasthttp

import (
	"github.com/edunx/lua"
	pub "github.com/edunx/rock-public-go"
)

const (
	SERVERMT string = "ROCK_FASTHTTP_SERVER_GO_MT"
	ACCESSLOG string = "http_time,server_addr,server_port,remote_addr,host,path"
)

func LuaInjectServerApi(L *lua.LState , parent *lua.LTable) {
	mt := L.NewTypeMetatable( SERVERMT )
	L.SetField(mt , "__index" , L.NewFunction(serverIndex))
	L.SetField(mt , "__newindex" , L.NewFunction(serverNewindex))

	L.SetField(parent , "server" , L.NewFunction(CreateServerUserData))
}

func CreateServerUserData(L *lua.LState) int {
	opt := L.CheckTable(1)

	v := &Server{
		C: Config{
			listen:    opt.CheckSocket("listen" , L),
			protocol:  opt.CheckString("protocol" , "tcp"),
			routers:   opt.CheckString("router" , "router"),
			handler:   opt.CheckString("handler" , "handler"),
			unknown:   opt.CheckString("default" , "default"),
			reuseport: opt.CheckString("reuseport" , "off"),
			keepalive: opt.CheckString("keepalive" , "on"),
			daemon:    opt.CheckString("daemon" , "off"),
			accessLog: opt.CheckString("access_log" , ACCESSLOG),
			accessFormat: opt.CheckString("access_format" , "json"),
			accessRegion: opt.CheckString("access_region" , "x-real-ip"),
		},

		region: CheckRegionUserData(L , opt.RawGetString("region")),
		access: CheckTunnelUserData(L , opt.RawGetString("access")),
	}

	ud := L.NewUserDataByInterface(v , SERVERMT)
	L.Push(ud)
	return 1
}

func serverIndex(L *lua.LState) int {
	self := CheckServerUserData(L , 1)
	name := L.CheckString(2)
	pub.Out.Err("Get name : %s" , name)
	switch name {
	case "start":
		L.Push(L.NewFunction(self.startByLua))
	}

	return 1
}

func serverNewindex(L *lua.LState) int {
	return 0
}

func (self *Server) startByLua(L *lua.LState) int {

	if e := self.Start() ; e != nil {
		L.Push(lua.LString( e.Error() ))
		pub.Out.Err("fashttp server start fail , err: %v" , e)
		return 1
	}

	L.Push(lua.LNil)
	return 1
}


func(self *Server) ToUserData(L *lua.LState) *lua.LUserData {
	return L.NewUserDataByInterface( self , SERVERMT)
}

