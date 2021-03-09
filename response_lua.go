package fasthttp

import (
	"github.com/edunx/lua"
	pub "github.com/edunx/rock-public-go"
)

func injectResponseApi(L *lua.LState, parent *lua.LTable) {

	respTab := L.CreateTable(0 , 3)

	L.SetField(respTab , "set_body" ,   L.NewFunction( setBody ))
	L.SetField(respTab , "set_code" ,   L.NewFunction( setCode ))
	L.SetField(respTab , "set_header" , L.NewFunction( setheader ))

	L.SetField(parent , "response" , respTab)
}

func setBody(L *lua.LState) int {
	ctx := CheckRequestCtx( L )
	body := L.CheckString(1)
	ctx.Response.SetBody( pub.S2B( body ))
	return 0
}

func setCode(L *lua.LState) int {
	code := L.CheckInt(1)
	ctx := CheckRequestCtx( L )
	ctx.Response.SetStatusCode( code )
	return 0
}

func setheader(L *lua.LState) int {
	tab := L.CheckTable( 1 )
	kvs := CheckKeyValUserDatUserDataSlice(L , tab)
	ctx := CheckRequestCtx( L )
	size := len(kvs)
	if size == 0 {
		return 0
	}

	var item *KeyVal
	for i:=0 ; i < size ;i++ {
		item = kvs[i]
		ctx.Response.Header.Set(item.Key , item.Val)
	}
	return  0
}

