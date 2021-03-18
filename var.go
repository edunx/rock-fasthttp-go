package fasthttp

import (
	pub "github.com/edunx/rock-public-go"
	"github.com/valyala/fasthttp"
	"net"
)

func httpServerPort(ctx *fasthttp.RequestCtx) int {
	addr := ctx.LocalAddr()
	x, ok := addr.(*net.TCPAddr)
	if !ok {
		return 0
	}
	return x.Port
}

func httpRemotePort(ctx *fasthttp.RequestCtx) int {
	addr := ctx.RemoteAddr()
	x, ok := addr.(*net.TCPAddr)
	if !ok {
		return 0
	}
	return x.Port
}

func httpRegionCityId(ctx *fasthttp.RequestCtx) int {
	obj := ctx.UserValue("region_cityid")
	if val , ok := obj.(int); ok {
		return val
	}
	return 0
}

func httpRegionInfoRaw(ctx *fasthttp.RequestCtx) []byte {
	obj := ctx.UserValue("region_info")
	if val , ok := obj.([]byte); ok {
		return val
	}
	return pub.S2B("")
}