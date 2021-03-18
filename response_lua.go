package fasthttp

import (
	"github.com/edunx/lua"
	pub "github.com/edunx/rock-public-go"
	base "github.com/edunx/rock-base-go"
	"strings"
)

func injectResponseApi(L *lua.LState, parent *lua.LTable) {

	respTab := L.CreateTable(0 , 3)

	L.SetField(respTab , "say" ,   L.NewFunction( responseSay  ))
	L.SetField(respTab , "append" ,   L.NewFunction( responseAppend  ))
	L.SetField(respTab , "exit" ,   L.NewFunction( responseExit ))
	L.SetField(respTab , "header" , L.NewFunction( responseHeader ))

	L.SetField(parent , "response" , respTab)
}

func responseAppend(L *lua.LState) int {
	ctx := CheckRequestCtx( L )
	body := L.CheckString(1)
	ctx.Response.AppendBody( pub.S2B( body ))
	return 0
}

func responseSay(L *lua.LState) int {
	ctx := CheckRequestCtx( L )
	n := L.GetTop()
	data := make([]string , n)

	for i:=1 ;i<=n ;i++ {
		data[i - 1]	 = L.CheckString(i)
	}
	ctx.Response.SetBody( pub.S2B( strings.Join( data , "")))
	return 0
}

func responseExit(L *lua.LState) int {
	code := L.CheckInt(1)
	ctx := CheckRequestCtx( L )
	ctx.Response.SetStatusCode( code )
	return 0
}

func responseHeader(L *lua.LState) int {
	tab := L.CheckTable( 1 )
	kvs := base.CheckKeyValUserDatUserDataSlice(L , tab)
	ctx := CheckRequestCtx( L )
	size := len(kvs)
	if size == 0 {
		return 0
	}

	var item *base.KeyVal
	for i:=0 ; i < size ;i++ {
		item = kvs[i]
		ctx.Response.Header.Set(item.Key , item.Val)
	}
	return  0
}

