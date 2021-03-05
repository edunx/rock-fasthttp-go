package fasthttp

import (
	"github.com/edunx/lua"
)

func LuaInjectApi(L *lua.LState , parent *lua.LTable) {
	fasthttpTab := L.CreateTable(0 , 2)
	LuaInjectServerApi(L , fasthttpTab)
	L.SetField(parent , "fasthttp" , fasthttpTab)
}
