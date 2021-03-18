package fasthttp

import (
	"github.com/edunx/lua"
)


func injectLoggerApi(L *lua.LState , parent *lua.LTable ) {
	loggerTab := L.CreateTable(0 , 1)
	L.SetField(loggerTab , "json" , L.NewFunction( newLoggerJson ))

	L.SetField(parent , "logger" , loggerTab )
}

func newLoggerJson( L *lua.LState ) int {
	vlog := new(vlogger)

	val := L.CheckString(1)
	vlog.New(val)
	vlog.encode = vlog.Json

	L.ExData.Set("logger" , vlog)

	return 0
}