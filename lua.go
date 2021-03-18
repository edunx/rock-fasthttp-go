package fasthttp

import (
	"github.com/edunx/lua"
)

func LuaInjectApi(L *lua.LState , parent *lua.LTable) {
	tab := L.CreateTable(0 , 2)
	LuaInjectServerApi(L , tab)
	L.SetField(parent , "fasthttp" , tab)
}