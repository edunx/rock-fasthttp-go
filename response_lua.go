package fasthttp

import (
	"github.com/edunx/lua"
	pub "github.com/edunx/rock-public-go"
	"github.com/valyala/fasthttp"
)

func injectResponseApi(L *lua.LState, parent *lua.LTable) {
	respTab := L.CreateTable(0 , 2)

	L.SetField(respTab , "set_body" , L.NewFunction( setBody ))
	L.SetField(parent , "response" , respTab)
}

func setBody(L *lua.LState) int {
	ctx := L.GetExdata().(*fasthttp.RequestCtx )
	body := L.CheckString(1)
	ctx.Response.SetBody( pub.S2B( body ))
	return 0
}

