package fasthttp

import (
	"fmt"
	base "github.com/edunx/rock-base-go"
	pub "github.com/edunx/rock-public-go"
	"github.com/valyala/fasthttp"
	"runtime/debug"
)

func panicHandler( ctx *fasthttp.RequestCtx , val interface{} ) {
	ctx.Response.SetStatusCode( 500)
	e := fmt.Sprintf("%v %s" , val , debug.Stack() )
	pub.Out.Err(e)
	ctx.Response.SetBodyString( e )
}

func (v *vHandler) SetHeader(ctx *fasthttp.RequestCtx) {
	if v.header == nil {
		return
	}

	size := len(v.header)
	if size == 0 {
		return
	}

	var item *base.KeyVal
	for i:=0 ; i < size ;i++ {
		item = v.header[i]
		ctx.Response.Header.Set(item.Key , item.Val)
	}
}

func (v *vHandler) SetBody( ctx *fasthttp.RequestCtx) {
	ctx.Response.SetStatusCode( v.code )
	ctx.Response.SetBodyString( v.body )
}

func (v *vHandler) Do( ctx *fasthttp.RequestCtx ) {
	v.SetHeader( ctx )
	v.SetBody( ctx )

	call(ctx , v.hook)
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

func (vhc *vHandlerChains) Do( ctx *fasthttp.RequestCtx ) {
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

		vh.Do( ctx )
		if vh.eof == "on" {
			return //结束匹配
		}
	}
}