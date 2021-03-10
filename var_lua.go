package fasthttp

import (
	"github.com/edunx/lua"
	pub "github.com/edunx/rock-public-go"
)

const (
	VARMT = "ROCK_FASTHTTP_VAR_GO_MT"
)

func injectVarApi(L *lua.LState , parent *lua.LTable) {
	varMT := L.NewTypeMetatable( VARMT )

	L.SetField(varMT , "__index" , L.NewFunction( varIndex ))
	L.SetField(varMT , "__newindex" , L.NewFunction( varNewIndex ))

	L.SetField(parent , "var" , L.NewUserDataByInterface(nil , VARMT )) //默认
}

func varIndex(L *lua.LState) int {
	ctx := CheckRequestCtx( L )
	name := L.CheckString(2)

	switch name {
	case "remote_addr":
		L.Push(lua.LString(ctx.RemoteAddr().String()))
	case "host":
		L.Push(lua.LString(ctx.Host()))
	case "path":
		L.Push(lua.LString(ctx.Request.URI().Path()))
	case "args":
		L.Push(lua.LString(ctx.QueryArgs().String()))
	case "full_uri":
		L.Push(lua.LString(ctx.Request.URI().FullURI()))
	case "request_uri":
		L.Push(lua.LString(ctx.Request.URI().RequestURI()))
	case "http_time":
		L.Push(lua.LNumber(ctx.Time().Unix()))
	case "cookie_raw":
		L.Push(lua.LString(ctx.Request.Header.Peek("cookie") ) )
	case "header_raw":
		L.Push(lua.LString(ctx.Request.Header.String() ) )
	case "content_length":
		L.Push(lua.LNumber(ctx.Request.Header.ContentLength()))
	case "content_type":
		L.Push(lua.LString(ctx.Request.Header.ContentType()) )
	case "content":
		L.Push(lua.LString(ctx.Request.Body()))
	default:
		size := len(name)
		if size > 4 && name[:4] == "arg_" {
			L.Push(lua.LString(ctx.QueryArgs().Peek( name[4:] )))
			return 1
		}

		if size > 5 && name[:5] == "post_" {
			L.Push(lua.LString(ctx.Request.PostArgs().Peek( name[5:])))
			return 1
		}

		if size > 5 && name[:5] == "http_" {
			L.Push( lua.LString(ctx.Request.Header.Peek(name[5:])) )
			return 1
		}

		if size > 7 && name[:7] == "cookie_" {
			L.Push(lua.LString(ctx.Request.Header.Cookie( name[7:])))
			return 1
		}

		if size > 6 && name[:6] == "param_" {
			obj := ctx.UserValue( name[6:])
			if val , ok := obj.(string); ok {
				L.Push(lua.LString( val ))
				return  1
			}

			L.Push(lua.LNil)
			pub.Out.Err(" param not found ")

			return 1
		}
	}

	return 1
}

func varNewIndex(L *lua.LState) int {
	return 0
}