package fasthttp

import (
	"github.com/edunx/lua"
	"github.com/fasthttp/router"
	"strings"
)

func newState() *lua.LState {
	vm := lua.NewState( )
	r  := router.New()

	vm.SetExdata(r)
	tab := vm.CreateTable( 0 , 1)
	injectHttpFuncsApi(vm , tab)
	vm.SetGlobal("http" , tab)

	return  vm
}

func CheckRegionUserData( L *lua.LState , v lua.LValue) region {

	ud , ok := v.(*lua.LUserData)
	if !ok {
		L.RaiseError("region must be userdata , go %T" , v)
		return nil
	}

	r , ok := ud.Value.(region)
	if !ok {
		L.RaiseError("region must have search , but not found")
		return nil
	}

	return r
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

func CheckHandlers( L *lua.LState) ([]string , int ) {
	data := L.CheckString( 2 )
	val  := strings.Split( data , ",")
	return val , len(val)
}

//func CheckHandlers( L *lua.LState ) ([]*vHandler , int) {
//	n := L.GetTop()
//	pub.Out.Err("got top == %d" , n)
//	if n < 2 {
//		L.RaiseError("not found handler fail")
//		return nil , 0
//	}
//
//	rc := make([]*vHandler , n - 1)
//
//	for i := 2 ; i <= n ; i++ {
//		rc[ i - 2 ] = CheckHandler(L , i)
//	}
//
//	pub.Out.Debug("got h == %v , size == %d" , rc , n)
//
//	return rc , n - 1
//}

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
		if rlen != isize + 1 {
			continue
		}

		// *risk1, *risk2
		if item[0] == '*' && item[1:] == risk {
			return true
		}

		if item[isize] == '*' && item[:rlen] == risk {
			return true
		}
	}

	return false
}
