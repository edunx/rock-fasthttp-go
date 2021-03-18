package fasthttp

import (
	"github.com/edunx/lua"
	pub "github.com/edunx/rock-public-go"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func injectRouterApi(L *lua.LState , parent *lua.LTable) {

	routerMT := L.NewTypeMetatable( ROUTERMT )
	L.SetField(routerMT , "__index" , L.NewFunction( routerUserDataIndex ))
	L.SetField(parent , "router" , L.NewUserDataByInterface( CheckExDataRouter( L ) , ROUTERMT ))
}

func CheckExDataRouter( L *lua.LState) *router.Router {
	r  , ok := L.ExData.Get("router").(*router.Router)
	if ok {
		return r
	}

	panic(" expect invalid router")
	pub.Out.Err("expect invalid router")

	return nil
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
	case "region":
		L.Push(regionIndexFn(L))
	case "access_push_off":
		L.Push(accessPushOffIndexFn(L))
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

func regionIndexFn( L *lua.LState ) *lua.LFunction {
	fn := func(vm *lua.LState) int {

		//判断是不是子线程
		val := vm.CheckString(1)
		if vm.Parent == nil {
			L.ExData.Set("region" , val)
			return 0
		}

		if ctx := CheckRequestCtx( vm ); ctx != nil {
			ctx.SetUserValue("region" , val)
		}

		return 0
	}

	return L.NewFunction( fn )
}

func accessPushOffIndexFn(L *lua.LState) *lua.LFunction {
	fn := func(vm *lua.LState) int {
		if vm.Parent == nil {
			L.ExData.Set("access_push_off" , "off")
			return 0
		}

		if ctx := CheckRequestCtx(vm); ctx != nil {
			ctx.SetUserValue("access_push_off" , "off")
		}
		return 0

	}

	return L.NewFunction( fn )
}
