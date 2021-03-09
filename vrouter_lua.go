package fasthttp

import (
	"github.com/edunx/lua"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func injectRouterApi(L *lua.LState , parent *lua.LTable) {

	routerMT := L.NewTypeMetatable( ROUTERMT )
	L.SetField(routerMT , "__index" , L.NewFunction( routerUserDataIndex ))

	L.SetField(parent , "router" , L.NewUserDataByInterface( L.GetExdata().(*router.Router) , ROUTERMT ))

}

func routerUserDataIndex (L *lua.LState) int {
	r := CheckRouterUserData(L, 1)
	method := L.CheckString(2)

	switch method {

	case "GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE":
		L.Push(L.NewFunction(func(vm *lua.LState) int {
			path := vm.CheckString(1)
			chains := CheckHandlerChains(vm)
			r.Handle(method, path, func(ctx *fasthttp.RequestCtx) { chains.Do( L , ctx ) })
			return 0
		}))

	case "ANY":
		L.Push(L.NewFunction(func(vm *lua.LState) int {
			path := vm.CheckString(1)
			chains := CheckHandlerChains(vm)
			r.ANY(path, func(ctx *fasthttp.RequestCtx) { chains.Do( L , ctx ) })
			return 0
		}))

	case "not_found":
		L.Push(L.NewFunction(func(vm *lua.LState) int {
			chains := CheckHandlerChains(vm)
			r.NotFound = func(ctx *fasthttp.RequestCtx) { chains.Do(L , ctx ) }
			return 0
		}))
	}

	return 1
}
