package fasthttp

import "github.com/edunx/lua"

const (
	VRULEMT = "ROCK_VRULE_GO_MT"
)

func injectRuleApi(L *lua.LState , parent *lua.LTable) {
	L.SetField(parent , "rule" , L.NewFunction( createRuleLightUserData ))
}

func createRuleLightUserData(L *lua.LState) int {
	opt := L.CheckTable(1)

	val := vRule{
		name: opt.CheckString("var" , ""),
		method: opt.CheckString("method" , "eq"),
		value: opt.RawGetString("value"),
	}

	L.Push( L.NewLightUserData( val ))
	return 1
}

