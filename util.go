package fasthttp

import (
	"github.com/edunx/lua"
	pub "github.com/edunx/rock-public-go"
	tp "github.com/edunx/rock-transport-go"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func newThread(L *lua.LState ) *thread {
	co , fn := L.NewThread()
	co.Parent = L

	return &thread{ co , fn }
}

func newState() *lua.LState {

	//创建虚拟解释的虚拟机
	vm := lua.NewState( )
	r  := router.New()

	vm.ExData.Set("router" , r) //注册虚拟机 router 对象

	//vm.ExData.Set("logger" , func( ctx *fasthttp.RequestCtx) { ctx.Logger().Printf("default")})

	tab := vm.CreateTable( 0 , 3)
	injectHttpFuncsApi(vm , tab)
	vm.SetGlobal("http" , tab)

	r.PanicHandler = panicHandler
	return  vm
}

func call(ctx *fasthttp.RequestCtx , hook *lua.LFunction ) {
	if hook == nil {
		return
	}

	vctx := ctx.UserValue("vctx").(*vContext)
	if vctx.vth == nil {
		vctx.vth = vctx.vrr.Co.Get().(*thread)
	}

	vctx.vth.co.ExData.Set("ctx" , ctx)
	vctx.vth.co.Push( hook )
	if e := vctx.vth.co.PCall(0 , 0 , nil) ; e != nil {
		pub.Out.Err("http hook run err: %v", e)
	}

	vctx.vth.co.ExData.Set("ctx" , nil)
	vctx.vrr.Co.Put(vctx.vth)

}

func CheckRegionUserData( L *lua.LState , v lua.LValue) region {
	ud , ok := v.(*lua.LUserData)
	if !ok {
		//L.RaiseError("region must be userdata , go %T" , v)
		return nil
	}

	r , ok := ud.Value.(region)
	if !ok {
		//L.RaiseError("region must have search , but not found")
		return nil
	}


	return r
}

func CheckTunnelUserData(L *lua.LState , v lua.LValue) tp.Tunnel {
	ud , ok := v.(*lua.LUserData)
	if !ok {
		pub.Out.Err("access log tunnel got nil")
		return nil
	}

	obj , ok := ud.Value.(tp.Tunnel)
	if !ok {
		pub.Out.Err("access log tunnel got invalid")
		return nil
	}

	return obj

}

func CheckServerUserData(L *lua.LState , idx int ) *Server {

	ud := L.CheckUserData( idx )

	v  , ok := ud.Value.(*Server)
	if ok {
		return v
	}

	L.TypeError(idx , lua.LTUserData)
	return nil

}

func CheckRouterUserData(L *lua.LState , idx int ) *router.Router{

	ud := L.CheckUserData( idx )

	v  , ok := ud.Value.(*router.Router)
	if ok {
		return v
	}

	L.TypeError(idx , lua.LTUserData)
	return nil
}

func CheckLuaFunctionByTable(L *lua.LState , opt *lua.LTable , key string ) *lua.LFunction {
	v := opt.RawGetString(key)

	fn , ok := v.(*lua.LFunction)
	if !ok {
		//L.RaiseError("%s must be function , got fail" , key)
		return nil
	}

	return fn
}

func CompareRule( rule []string , risk string , rlen int)  bool {
	size := len(rule)
	if size == 0 {
		return false
	}

	var item string
	var isize int
	for i := 0 ;i < size ; i++ {
		item = rule[i]
		if item == "*" {
			return true
		}

		if item == risk {
			return true
		}

		isize = len( item )
		if rlen < isize + 1 {
			continue
		}

		// *risk1, *risk2
		if item[0] == '*' && item[1:] == risk[:isize - 1] {
			return true
		}

		if item[isize] == '*' && item[:rlen] == risk[rlen - isize: rlen - 1] {
			return true
		}
	}

	return false
}

func CheckRequestCtx(co *lua.LState) *fasthttp.RequestCtx {
	ctx , ok := co.ExData.Get("ctx").(*fasthttp.RequestCtx)
	if ok {
		return ctx
	}

	return nil
}
