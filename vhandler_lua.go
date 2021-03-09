package fasthttp

import (
	"github.com/edunx/lua"
	"strings"
)

func injectHandlerApi(L *lua.LState , parent *lua.LTable) {
	L.SetField(parent , "handler" , L.NewFunction( CreateHandlerUserData ))
}

func CheckHandler(L *lua.LState , idx int) *vHandler {

	ud := L.CheckLightUserData( idx )
	v , ok := ud.Value.(*vHandler)
	if !ok {
		L.RaiseError("#%d must be http.vhandler , got fail" , idx)
		return nil
	}

	return v
}

func CreateHandlerUserData(L *lua.LState) int {
	opt := L.CheckTable( 1 )

	v := &vHandler{
		rule:  strings.Split(opt.CheckString("rule" , "*") , ","),
		tag:   opt.CheckString("tag" , "null"),
		header: CheckKeyValUserDatUserDataSlice(L , opt.RawGetString("header")),
		body: opt.CheckString("body" , "null"),
		code: opt.CheckInt("code" , 400),
		eof: opt.CheckString("eof" , "on"),
		bodyEncode: opt.CheckString("body_encode" , ""),
		bodyEncodeMin: opt.CheckInt("body_encode_min" , 100),
		hook: CheckLuaFunctionByTable(L , opt , "hook"),
	}

	L.Push(L.NewLightUserData( v ))

	return 1
}

func CheckHandlerChains( L *lua.LState )  vHandlerChains {
	n := L.GetTop()

	vhc := vHandlerChains{ cap:0 }
	if n < 2 {
		L.RaiseError("not found handler fail")
		return vhc
	}

	vhc.cap = n - 1
	vhc.data = make([]interface{} , n - 1 )
	vhc.mask = make([]int , n - 1 )

	var val lua.LValue
	for i := 2; i <= n; i++ {
		val = L.Get( i )
		switch val.Type() {
		case lua.LTString:
			vhc.Store( val.String() , VHSTRING  , i - 2 )

		case lua.LTLightUserData:
			obj , ok := val.(*lua.LightUserData).Value.(*vHandler)
			if ok {
				vhc.Store( obj , VHANDLER , i - 2 )
			} else {
				L.RaiseError("#%d must be string or http.handler , got fail" , i)
				goto ERR
			}

		default:
			L.RaiseError("#%d must be string or http.handler , got fail" , i)
			goto ERR
		}
	}

	return vhc //正常返回

	ERR:
		return vHandlerChains{ cap: 0 }
}