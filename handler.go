package fasthttp

import (
	"github.com/edunx/lua"
	pub "github.com/edunx/rock-public-go"
	"github.com/valyala/fasthttp"
)

func handler( ctx *fasthttp.RequestCtx ) {
	r := cvr.load( pub.B2S( ctx.Host() ) )
	r.Handler(ctx)
}

func (v *vHandler) SetHeader(ctx *fasthttp.RequestCtx) {
	if v.header == nil {
		return
	}

	size := len(v.header)
	if size == 0 {
		return
	}

	var item *KeyVal
	for i:=0 ; i < size ;i++ {
		item = v.header[i]
		ctx.Response.Header.Set(item.Key , item.Val)
	}
}

func (v *vHandler) SetBody( ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode( v.code )
	ctx.Response.SetBodyString( v.body )
}

func (v *vHandler) Call( L *lua.LState , ctx *fasthttp.RequestCtx ) {
	if v.hook == nil {
		return
	}

	th := v.Pool.Get().(*thread)
	th.co.Push( v.hook )
	th.co.SetExdata( ctx )

	if e := th.co.PCall(0 , 0 , nil) ; e != nil {
		pub.Out.Err("http hook run err: %v", e)
	}

	th.co.SetExdata( nil )
	v.Pool.Put( th )
}

func (v *vHandler) Set(L *lua.LState , ctx *fasthttp.RequestCtx ) {
	v.SetHeader( ctx )
	v.SetBody( ctx )
	v.Call( L , ctx )
}

func handlerLoop( ctx *fasthttp.RequestCtx , vhs []string , size int , L *lua.LState) {
	if size <= 0 {
		ctx.Response.SetStatusCode(400)
		ctx.Response.SetBodyString("not found handler")
		return
	}

	data := ctx.Request.Header.Peek( "risk ")
	rSize := len(data)
	risk := pub.B2S( data )

	var vh *vHandler
	for i := 0 ; i < size ;i++ {
		vh = cvm.load( vhs[i] )
		if !CompareRule( vh.rule , risk , rSize) {
			continue
		}

		vh.Set( L , ctx )
		if vh.eof == "on" {
			return //结束匹配
		}

	}

}

//func DoHandlerLoop( ctx *fasthttp.RequestCtx , vhs []*vHandler , size int , L *lua.LState) {
//	if size <= 0 {
//		ctx.Response.SetStatusCode(400)
//		ctx.Response.SetBodyString("not found handler")
//		return
//	}
//
//	data := ctx.Request.Header.Peek( "risk")
//	rSize := len(data)
//	risk:= pub.B2S( data )
//
//	var vh *vHandler
//	for i := 0 ; i < size ; i++ {
//		vh = vhs[i]
//		if !CompareRule( vh.rule , risk , rSize) {
//			continue
//		}
//
//		vh.Set( L , ctx )
//		if vh.eof == "on" {
//			return //结束匹配
//		}
//	}
//}
