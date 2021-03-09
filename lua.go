package fasthttp

import (
	"github.com/edunx/lua"
	"github.com/fasthttp/router"
)

func newHttpThreadState() *lua.LState {

	vm := lua.NewState( )
	r  := router.New()
	vm.SetExdata(r)
	tab := vm.CreateTable( 0 , 3)
	injectHttpFuncsApi(vm , tab)
	vm.SetGlobal("http" , tab)
	return  vm

}

func LuaInjectApi(L *lua.LState , parent *lua.LTable) {
	tab := L.CreateTable(0 , 2)
	LuaInjectServerApi(L , tab)
	L.SetField(parent , "fasthttp" , tab)
}
