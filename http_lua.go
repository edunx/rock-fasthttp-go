package fasthttp

import (
	"github.com/edunx/lua"
	base "github.com/edunx/rock-base-go"
)

const (
	ROUTERMT string = "ROCK_FASTHTTP_ROUTER_GO_MT"
)

func injectHttpFuncsApi(L *lua.LState , parent *lua.LTable) {

	L.SetField(parent , "keyval" , L.NewFunction( CreateKeyValUserData ))

	injectVarApi(L , parent)
	injectRuleApi(L , parent)
	injectRouterApi(L , parent)
	injectHandlerApi(L , parent)
	injectResponseApi(L , parent)

	base.LuaInjectApi(L , parent)
}

func CreateKeyValUserData(L *lua.LState) int {
	key := L.CheckString(1)
	val := L.CheckString(2)
	ud := L.NewLightUserData( &KeyVal{key , val } )
	L.Push(ud)
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