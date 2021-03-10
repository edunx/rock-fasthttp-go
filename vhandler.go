package fasthttp

import (
	"fmt"
	"github.com/edunx/lua"
	pub "github.com/edunx/rock-public-go"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func panicHandler( ctx *fasthttp.RequestCtx , val interface{} ) {
	ctx.Response.SetStatusCode( 500)
	ctx.Response.SetBodyString(fmt.Sprintf("%v" , val))
}

func handler( ctx *fasthttp.RequestCtx ) {
	vrr := cvr.load( pub.B2S( ctx.Host() ) )
	if vrr == nil {
		ctx.Response.SetStatusCode(500)
		ctx.Response.SetBody(pub.S2B("not found router"))
		return
	}

	r , ok := vrr.L.GetExdata().(*router.Router)
	if !ok {
		ctx.Response.SetStatusCode(500)
		ctx.Response.SetBody(pub.S2B("expect invalid router"))
		return
	}

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

	r := ctx.UserValue("router").(*vRouter)

	th := r.Co.Get().(*thread)
	th.co.Push( v.hook )
	th.co.SetExdata( ctx )

	if e := th.co.PCall(0 , 0 , nil) ; e != nil {
		pub.Out.Err("http hook run err: %v", e)
	}
	th.co.SetExdata( nil )
	r.Co.Put( th )
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
		if vh == nil {
			continue
		}

		if !CompareRule( vh.rule , risk , rSize) {
			continue
		}

		vh.Set( L , ctx )
		if vh.eof == "on" {
			return //结束匹配
		}

	}
}

func (vhc *vHandlerChains) Store( v interface{} , mask int  , cap int ) {
	if cap > vhc.cap {
		pub.Out.Err("vhandler overflower ,cap: %d , got: %d" , vhc.cap , cap)
		return
	}

	vhc.data[cap] = v
	vhc.mask[cap] = mask
}

func (vhc *vHandlerChains) notFound( ctx *fasthttp.RequestCtx ) {
	ctx.Response.SetStatusCode(404)
	ctx.Response.SetBodyString(ctx.Request.String() + " not found handler")
}

func (vhc *vHandlerChains) Do( L *lua.LState , ctx *fasthttp.RequestCtx ) {
	if vhc.cap == 0 {
		vhc.notFound( ctx )
		return
	}

	data := ctx.Request.Header.Peek( "risk")
	rSize := len(data)
	risk:= pub.B2S( data )

	var vh *vHandler
	for i := 0 ; i < vhc.cap ; i++ {

		switch vhc.mask[i] {
		case VHSTRING:
			vh = cvm.load( vhc.data[i].(string))
		case VHANDLER:
			vh = vhc.data[i].(*vHandler)
		default:
			continue
		}

		if !CompareRule( vh.rule , risk , rSize) {
			continue
		}

		vh.Set( L , ctx )
		if vh.eof == "on" {
			return //结束匹配
		}
	}
}