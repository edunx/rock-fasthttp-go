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

func varIndex(co *lua.LState) int {
	ctx := CheckRequestCtx(co)
	name := co.CheckString(2)

	switch name {
	case "remote_addr":
		co.Push(lua.LString(ctx.RemoteIP().String()))
	case "remote_port":
		co.Push(lua.LNumber(httpRemotePort(ctx)))
	case "server_addr":
		co.Push(lua.LString(ctx.LocalIP().String()))
	case "server_port":
		co.Push(lua.LNumber(httpServerPort(ctx)))
	case "host":
		co.Push(lua.LString(ctx.Host()))
	case "uri":
		co.Push(lua.LString(ctx.Request.URI().Path()))
	case "args":
		co.Push(lua.LString(ctx.QueryArgs().String()))
	case "full_uri":
		co.Push(lua.LString(ctx.Request.URI().FullURI()))
	case "request_uri":
		co.Push(lua.LString(ctx.Request.URI().RequestURI()))
	case "http_time":
		co.Push(lua.LNumber(ctx.Time().Unix()))
	case "cookie_raw":
		co.Push(lua.LString(ctx.Request.Header.Peek("cookie") ) )
	case "header_raw":
		co.Push(lua.LString(ctx.Request.Header.String() ) )
	case "content_length":
		co.Push(lua.LNumber(ctx.Request.Header.ContentLength()))
	case "content_type":
		co.Push(lua.LString(ctx.Request.Header.ContentType()) )
	case "content":
		co.Push(lua.LString(ctx.Request.Body()))
	case "region_cityid":
		co.Push(lua.LNumber(httpRegionCityId( ctx )))
	case "region_info":
		co.Push(lua.LString(httpRegionInfoRaw( ctx )))
	case "ua":
		co.Push(lua.LString(ctx.Request.Header.UserAgent()))
	case "referer":
		co.Push(lua.LString(ctx.Request.Header.Referer()))
	case "status":
		co.Push(lua.LString(ctx.Response.StatusCode()))
	case "sent":
		co.Push(lua.LNumber(ctx.Response.Header.ContentLength()))

	default:
		size := len(name)
		if size > 4 && name[:4] == "arg_" {
			co.Push(lua.LString(ctx.QueryArgs().Peek( name[4:] )))
			return 1
		}

		if size > 5 && name[:5] == "post_" {
			co.Push(lua.LString(ctx.Request.PostArgs().Peek( name[5:])))
			return 1
		}

		if size > 5 && name[:5] == "http_" {
			co.Push( lua.LString(ctx.Request.Header.Peek(name[5:])) )
			return 1
		}

		if size > 7 && name[:7] == "cookie_" {
			co.Push(lua.LString(ctx.Request.Header.Cookie( name[7:])))
			return 1
		}

		if size > 6 && name[:6] == "param_" {
			obj := ctx.UserValue( name[6:])
			if val , ok := obj.(string); ok {
				co.Push(lua.LString( val ))
				return  1
			}

			co.Push(lua.LNil)
			pub.Out.Err(" param not found ")

			return 1
		}
	}

	return 1
}

func varNewIndex(co *lua.LState) int {
	ctx := CheckRequestCtx(co)
	name := co.CheckString(2)

	switch name  {
	case "uri":
		path := co.CheckString( 3 )
		ctx.Request.URI().SetPath( path )

	}

	return 0
}