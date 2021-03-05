package fasthttp

import (
	"github.com/edunx/lua"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"strings"
	"sync"
)

const (
	ROUTERMT string = "ROCK_FASTHTTP_ROUTER_GO_MT"
)

func injectHttpFuncsApi(L *lua.LState , parent *lua.LTable) {
	routerMT := L.NewTypeMetatable( ROUTERMT )
	L.SetField(routerMT , "__index" , L.NewFunction( routerUserDataIndex ))

	L.SetField(parent , "keyval" , L.NewFunction( CreateKeyValUserData ))
	L.SetField(parent , "router" , L.NewUserDataByInterface( L.GetExdata().(*router.Router) , ROUTERMT ))
	L.SetField(parent , "handler" , L.NewFunction( CreateHandlerUserData ))

	injectResponseApi(L , parent)
}

func CreateKeyValUserData(L *lua.LState) int {
	key := L.CheckString(1)
	val := L.CheckString(2)
	ud := L.NewLightUserData( &KeyVal{key , val } )
	L.Push(ud)
	return 1
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

	v.Pool = &sync.Pool{
		New: func() interface{} {
			co , fn := L.NewThread()
			return &thread{ co , fn }
		},
	}

	L.Push(L.NewLightUserData( v ))
	return 1
}

func CheckKeyValUserData(L *lua.LState , idx int) *KeyVal {
	ud := L.CheckLightUserData( idx )
	kv , ok := ud.Value.(*KeyVal)
	if !ok {
		L.RaiseError("#%d must keyval ,got fail" , idx)
		return  nil
	}

	return kv
}

func CheckKeyValUserDatUserDataSlice(L *lua.LState , lv lua.LValue) []*KeyVal {
	tab, ok := lv.(*lua.LTable)
	if !ok {
		//L.RaiseError("header must table , got fail")
		return nil
	}

	i := 0
	rc := make([]*KeyVal, tab.Len())

	tab.ForEach(func(idx lua.LValue, v lua.LValue) {
		if idx.Type() != lua.LTNumber {
			L.RaiseError("header must arr , got fail")
			return
		}

		ud, ok := v.(*lua.LightUserData)
		if !ok {
			L.RaiseError("header must lightuserdata, got fail")
			return
		}

		val, ok := ud.Value.(*KeyVal)
		if !ok {
			L.RaiseError("header must keyval, got fail")
			return
		}

		rc[i] = val
		i++
	})

	return rc
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


func routerUserDataIndex (L *lua.LState) int {
	r := CheckRouterUserData(L, 1)
	method := L.CheckString(2)

	switch method {

	case "GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE":
		L.Push(L.NewFunction(func(vm *lua.LState) int {
			path := vm.CheckString(1)
			handlers, size := CheckHandlers(vm)
			r.Handle(method, path, func(ctx *fasthttp.RequestCtx) { handlerLoop( ctx, handlers, size , L ) })
			return 0
		}))

	case "ANY":
		L.Push(L.NewFunction(func(vm *lua.LState) int {
			path := vm.CheckString(1)
			handlers, size := CheckHandlers(vm)
			r.ANY(path, func(ctx *fasthttp.RequestCtx) { handlerLoop(ctx, handlers, size , L ) })
			return 0
		}))

	case "not_found":
		L.Push(L.NewFunction(func(vm *lua.LState) int {
			handlers, size := CheckHandlers(vm)
			r.NotFound = func(ctx *fasthttp.RequestCtx) { handlerLoop(ctx, handlers, size , L ) }
			return 0
		}))
	}

	return 1
}

