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
		L.Push(handleIndexFn(L , method , r ))
	case "ANY":
		L.Push(anyIndexFn(L , r))
	case "not_found":
		L.Push(notFoundIndexFn( L , r ))
	case "file":
		L.Push(fileIndexFn(L , r))
	}

	return 1
}

func handleIndexFn(L *lua.LState , method string , r *router.Router ) *lua.LFunction {
	fn := func(vm *lua.LState) int {
		path := vm.CheckString(1)
		chains := CheckHandlerChains(vm)
		r.Handle(method, path, func(ctx *fasthttp.RequestCtx) { chains.Do( ctx ) })
		return 0
	}
	return L.NewFunction( fn )
}

func anyIndexFn(L *lua.LState , r *router.Router) *lua.LFunction {
	fn := func(vm *lua.LState) int {
		path := vm.CheckString(1)
		chains := CheckHandlerChains(vm)
		r.ANY(path, func(ctx *fasthttp.RequestCtx) { chains.Do( ctx ) })
		return 0
	}

	return L.NewFunction(fn)
}

func notFoundIndexFn(L *lua.LState , r *router.Router) *lua.LFunction {
	fn := func(vm *lua.LState) int {
		chains := CheckHandlerChains(vm)
		r.NotFound = func(ctx *fasthttp.RequestCtx) { chains.Do( ctx ) }
		return 0
	}
	return L.NewFunction( fn )

}

func fileIndexFn( L *lua.LState , r *router.Router ) *lua.LFunction {
	fn := func(vm *lua.LState ) int {
		n := vm.GetTop()
		path := vm.CheckString(1)
		root := vm.CheckString(2)
		fs := &fasthttp.FS{
			Root: root,
			IndexNames: []string{"index.html"},
			GenerateIndexPages: true,
			AcceptByteRange: true,
		}

		if n == 3 {
			fn := vm.CheckFunction( 3 )
			fs.PathRewrite = func(ctx *fasthttp.RequestCtx) []byte {
				call(ctx , fn)
				return ctx.Path()
			}
		}

		r.ServeFilesCustom(path , fs)

		return 0
	}

	return L.NewFunction(fn)
}