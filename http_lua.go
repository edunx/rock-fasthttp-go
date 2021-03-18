package fasthttp

import (
	"github.com/edunx/lua"
	base "github.com/edunx/rock-base-go"
)

const (
	ROUTERMT string = "ROCK_FASTHTTP_ROUTER_GO_MT"
)

func injectHttpFuncsApi(L *lua.LState , parent *lua.LTable) {

	injectVarApi(L , parent)
	injectRuleApi(L , parent)
	injectRouterApi(L , parent)
	injectHandlerApi(L , parent)
	injectResponseApi(L , parent)
	injectLoggerApi(L , parent)

	base.LuaInjectApi(L , parent)
}